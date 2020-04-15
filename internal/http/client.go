// Package http send request to get public ip and response to dyndns provider
package http

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/golgoth31/homedynip/internal/dns"

	jsoniter "github.com/json-iterator/go"
)

const goodExitCode = 0

// NewClient returns a new homedynip client
func NewClient() *Client {
	return &Client{}
}

func (c *Client) client() *http.Client {
	transp := http.DefaultTransport.(*http.Transport).Clone()
	transp.TLSClientConfig.InsecureSkipVerify = c.Config.GetBool("client.insecure")

	client := http.DefaultClient
	client.Transport = transp

	return client
}

func (c *Client) url() (string, string) {
	c.Log.Info().Msgf("service: %s", c.Config.GetString("client.service"))

	var serviceURL string

	var respField string

	switch c.Config.GetString("client.service") {
	case "ipify":
		serviceURL = "https://api.ipify.org?format=json"
		respField = "ip"
	case "httpbin":
		serviceURL = "https://httpbin.org/ip"
		respField = "origin"
	case "custom":
		serviceURL = c.Config.GetString("client.url")
		respField = "ip"
	default:
		serviceURL = "https://api.ipify.org?format=json"
		respField = "ip"
	}

	c.Log.Info().Msgf("url: %s", serviceURL)

	return serviceURL, respField
}

func (c *Client) getIP() (string, error) {
	client := c.client()
	url, field := c.url()

	resp, err := client.Get(url)
	if err != nil {
		c.Log.Info().Msgf("could not get answer: %v", err)
		return "", err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.Log.Info().Msg("unable to close body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.Log.Info().Msgf("could not read body: %v", err)
		return "", err
	}

	output := jsoniter.Get(body, field).ToString()
	c.IP = &net.IPAddr{
		IP: net.ParseIP(output),
	}

	if c.Config.GetBool("client.showip") {
		fmt.Println(output)
		os.Exit(goodExitCode)
	}

	return output, nil
}

// GetIP returns the IP given by server
func (c *Client) GetIP() (string, error) {
	return c.getIP()
}

// WriteDNS writes IP to dyndns provider
func (c *Client) WriteDNS() error {
	dnsProvider, err := dns.New(c.Config, c.IP.String(), c.Log)
	if err != nil {
		c.Log.Info().Msgf("Unknown DNS provider: %v", err)
		return err
	}

	if err := dnsProvider.Write(); err != nil {
		c.Log.Info().Msgf("Unable to write into DNS provider: %v", err)
		return err
	}

	return nil
}

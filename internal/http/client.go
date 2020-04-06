// Package http send request to get public ip and response to dyndns provider
package http

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/golgoth31/homedynip/internal/dns/ovh"
	jsoniter "github.com/json-iterator/go"
)

// NewClient returns a new homedynip client
func NewClient() *Client {
	return &Client{}
}

func (c *Client) client() *http.Client {
	transp := http.DefaultTransport.(*http.Transport).Clone()
	transp.TLSClientConfig.InsecureSkipVerify = c.Config.GetBool("insecure")

	client := http.DefaultClient
	client.Transport = transp

	return client
}

func (c *Client) url() (string, string) {
	log.Printf("service: %s", c.Config.GetString("client.service"))

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

	log.Printf("url: %s", serviceURL)

	return serviceURL, respField
}

func (c *Client) getIP() (string, error) {
	client := c.client()
	url, field := c.url()

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("could not get answer: %v", err)
		return "", err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Print("unable to close body")
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not read body: %v", err)
		return "", err
	}

	output := jsoniter.Get(body, field).ToString()
	c.IP = &net.IPAddr{
		IP: net.ParseIP(output),
	}

	return output, nil
}

// GetIP returns the IP given by server
func (c *Client) GetIP() (string, error) {
	return c.getIP()
}

// WriteDNS writes IP to dyndns provider
func (c *Client) WriteDNS() error {
	switch c.Config.GetString("client.dns") {
	case "ovh":
		dnsClient := &ovh.Ovh{
			Username: c.Config.GetString("ovh.username"),
			Password: c.Config.GetString("ovh.password"),
			Hostname: c.Config.GetString("ovh.hostname"),
			IP:       c.IP.String(),
		}
		if err := dnsClient.Write(); err != nil {
			log.Printf("Unable to write into DNS provider: %v", err)
			return err
		}
	default:
		log.Printf("Unknown DNS driver: %s", c.Config.GetString("client.dns"))
		return fmt.Errorf("unknow DNS driver: %s", c.Config.GetString("client.dns"))
	}

	return nil
}

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
	var serviceUrl string
	var respField string
	switch c.Config.GetString("client.service") {
	case "ipify":
		serviceUrl = "https://api.ipify.org?format=json"
		respField = "ip"
	case "httpbin":
		serviceUrl = "https://httpbin.org/ip"
		respField = "origin"
	case "custom":
		serviceUrl = c.Config.GetString("client.url")
		respField = "ip"
	default:
		serviceUrl = "https://api.ipify.org?format=json"
		respField = "ip"
	}
	log.Printf("url: %s", serviceUrl)
	return serviceUrl, respField
}

func (c *Client) getIp() (string, error) {
	client := c.client()
	url, field := c.url()
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("could not get answer: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("could not read body: %v", err)
		return "", err
	}
	output := jsoniter.Get(body, field).ToString()
	c.Ip = &net.IPAddr{
		IP: net.ParseIP(output),
	}
	return output, nil
}
func (c *Client) GetIp() (string, error) {
	return c.getIp()
}

func (c *Client) WriteDNS() error {
	switch c.Config.GetString("client.dns") {
	case "ovh":
		dnsClient := &ovh.Ovh{
			Username: c.Config.GetString("ovh.username"),
			Password: c.Config.GetString("ovh.password"),
			Hostname: c.Config.GetString("ovh.hostname"),
			Ip:       c.Ip.String(),
		}
		if err := dnsClient.Write(); err != nil {
			log.Printf("Unable to write into DNS provider: %v", err)
			return err
		}
	default:
		log.Printf("Unknown DNS driver: %s", c.Config.GetString("client.dns"))
		return fmt.Errorf("Unknow DNS driver: %s", c.Config.GetString("client.dns"))
	}
	return nil
}

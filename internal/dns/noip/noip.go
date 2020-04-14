// Package noip writes IP to dynhost
package noip

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const httpStatusCodeOK = 200

// Write function writes IP to dynhost
func (o *Noip) Write() error {
	client := http.DefaultClient
	url := fmt.Sprintf(
		"https://dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s",
		o.Hostname,
		o.IP,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(o.Username, o.Password)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			o.Log.Info().Msgf("unable to close body")
		}
	}()

	if resp.StatusCode != httpStatusCodeOK {
		o.Log.Error().Msgf("Error in response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	returnedBody := strings.SplitN(string(body), " ", 2)

	switch returnedBody[0] {
	case "nochg":
		o.Log.Info().Msgf("IP '%s' not changed", o.IP)
	case "good":
		o.Log.Info().Msgf("IP '%s' modified successfully", o.IP)
	default:
		o.Log.Error().Msgf("Unknown state: '%s'", returnedBody[0])
	}

	return nil
}

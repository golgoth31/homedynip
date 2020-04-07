// Package ovh writes IP to dynhost
package ovh

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Write function writes IP to dynhost
func (o *Ovh) Write() error {
	client := http.DefaultClient
	url := fmt.Sprintf("https://www.ovh.com/nic/update?system=dyndns&hostname=%s&myip=%s", o.Hostname, o.IP)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	o.Log.Info().Msg(string(body))

	return nil
}

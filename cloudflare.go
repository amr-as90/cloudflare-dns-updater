package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gregdel/pushover"
)

type cfConfig struct {
	zoneID           string
	dnsRecordID      []string
	dnsRecordNames   []string
	eMail            string
	apiKey           string
	pushoverAppToken string
	pushoverUserKey  string
}

type requestStruct struct {
	Comment    string `json:"comment"`
	Content    string `json:"content"`
	Name       string `json:"name"`
	Proxied    bool   `json:"proxied"`
	Ttl        int    `json:"ttl"`
	RecordType string `json:"type"`
}

type cfResponse struct {
	Success bool `json:"success"`
}

func (cfg cfConfig) cloudFlareUpdate(newIP string) error {

	for i := range cfg.dnsRecordID {

		reqStruct := requestStruct{
			Comment:    "Updated automatically via Go CloudFlare updater",
			Content:    newIP,
			Name:       cfg.dnsRecordNames[i],
			Proxied:    false,
			Ttl:        3600,
			RecordType: "A",
		}

		jsonData, err := json.Marshal(reqStruct)
		if err != nil {
			return err
		}

		cfURL := "https://api.cloudflare.com/client/v4/zones/" + cfg.zoneID + "/dns_records/" + cfg.dnsRecordID[i]

		req, err := http.NewRequest("PUT", cfURL, bytes.NewBuffer([]byte(jsonData)))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Email", cfg.eMail)
		req.Header.Set("X-Auth-Key", cfg.apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cfResponse := cfResponse{}

		err = json.Unmarshal(body, &cfResponse)
		if err != nil {
			return err
		}

		if cfResponse.Success {
			fmt.Printf("IP updated successfuly for: %s, new IP is: %s\n", cfg.dnsRecordNames[i], newIP)
			if cfg.pushoverAppToken != "" && cfg.pushoverUserKey != "" {
				app := pushover.New(cfg.pushoverAppToken)
				recipient := pushover.NewRecipient(cfg.pushoverUserKey)
				message := pushover.NewMessageWithTitle(fmt.Sprintf("IP of DNS record %s changed to %s", cfg.dnsRecordNames[i], newIP), "IP Changed")
				_, err := app.SendMessage(message, recipient)
				if err != nil {
					log.Printf("Error sending Pushover notification: %s", err)
				}
			}
		} else {
			fmt.Printf("Unable to update IP for %s. Something went wrong.\n", cfg.dnsRecordNames[i])
		}

	}
	return nil
}

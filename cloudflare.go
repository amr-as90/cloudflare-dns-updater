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
	records          map[string]string // name -> ID
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

// In cloudFlareUpdate method:
func (cfg cfConfig) cloudFlareUpdate(newIP string) error {
	if len(cfg.records) == 0 {
		return fmt.Errorf("no DNS records available")
	}

	for name, id := range cfg.records {
		if id == "" {
			return fmt.Errorf("empty DNS record ID for %s", name)
		}

		reqStruct := requestStruct{
			Comment:    "Updated automatically via Go CloudFlare updater",
			Content:    newIP,
			Name:       name, // Use the map key directly
			Proxied:    false,
			RecordType: "A",
		}

		jsonData, err := json.Marshal(reqStruct)
		if err != nil {
			return err
		}

		cfURL := "https://api.cloudflare.com/client/v4/zones/" + cfg.zoneID + "/dns_records/" + cfg.records[name] // Use the map value directly as the ID paramete

		req, err := http.NewRequest("PUT", cfURL, bytes.NewBuffer([]byte(jsonData)))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Auth-Email", cfg.eMail)
		req.Header.Set("Authorization", "Bearer "+cfg.apiKey)

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
			fmt.Printf("IP updated successfully for: %s, new IP is: %s\n", name, newIP)
			if cfg.pushoverAppToken != "" && cfg.pushoverUserKey != "" {
				app := pushover.New(cfg.pushoverAppToken)
				recipient := pushover.NewRecipient(cfg.pushoverUserKey)
				message := pushover.NewMessageWithTitle(
					fmt.Sprintf("IP of DNS record %s changed to %s", name, newIP),
					"IP Changed")
				_, err := app.SendMessage(message, recipient)
				if err != nil {
					log.Printf("Error sending Pushover notification: %s", err)
				}
			}
		} else {
			fmt.Printf("Unable to update IP for %s. Something went wrong.\n", name)
		}
	}
	return nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type cfConfig struct {
	zoneID      string
	dnsRecordID string
	eMail       string
	apiKey      string
}

type requestStruct struct {
	Comment    string `json:"comment"`
	Content    string `json:"content"`
	Name       string `json:"name"`
	Proxied    bool   `json:"proxied"`
	Ttl        int    `json:"ttl"`
	RecordType string `json:"type"`
}

func (cfg cfConfig) cloudFlareUpdate(newIP string) error {

	cfURL := "https://api.cloudflare.com/client/v4/zones/" + cfg.zoneID + "/dns_records/" + cfg.dnsRecordID

	reqStruct := requestStruct{
		Comment:    "Updated automatically via Go CloudFlare updater",
		Content:    newIP,
		Name:       "files.renderex.ae",
		Proxied:    false,
		Ttl:        3600,
		RecordType: "A",
	}

	json, err := json.Marshal(reqStruct)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", cfURL, bytes.NewBuffer([]byte(json)))
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

	fmt.Println(string(body))

	fmt.Printf("Response Status: %s", resp.Status)
	return nil
}

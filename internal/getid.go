package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cfg *CFConfig) GetDomainID() error {
	url := "https://api.cloudflare.com/client/v4/zones/" + cfg.ZoneID + "/dns_records"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Email", cfg.EMail)
	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)

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

	var dnsInfo struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
	}

	err = json.Unmarshal(body, &dnsInfo)
	if err != nil {
		return err
	}

	// Initialize records map if nil
	if cfg.Records == nil {
		cfg.Records = make(map[string]string)
	}

	// Find and store IDs for all configured record names
	for _, record := range dnsInfo.Result {
		if _, exists := cfg.Records[record.Name]; exists {
			cfg.Records[record.Name] = record.ID
		}
	}

	// Verify we found all records
	for name, id := range cfg.Records {
		if id == "" {
			return fmt.Errorf("could not find DNS record for %s", name)
		}
	}

	return nil
}

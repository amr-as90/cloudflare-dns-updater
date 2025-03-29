package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DNSinfo struct {
	Result []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Proxiable bool   `json:"proxiable"`
		Proxied   bool   `json:"proxied"`
		TTL       int    `json:"ttl"`
		Settings  struct {
		} `json:"settings,omitempty"`
		Meta struct {
		} `json:"meta"`
		Comment           string    `json:"comment"`
		Tags              []any     `json:"tags"`
		CreatedOn         time.Time `json:"created_on"`
		ModifiedOn        time.Time `json:"modified_on"`
		CommentModifiedOn time.Time `json:"comment_modified_on,omitempty"`
		Settings0         struct {
			FlattenCname bool `json:"flatten_cname"`
		} `json:"settings,omitempty"`
		Priority int `json:"priority,omitempty"`
	} `json:"result"`
	Success    bool  `json:"success"`
	Errors     []any `json:"errors"`
	Messages   []any `json:"messages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"result_info"`
}

func (cfg *cfConfig) GetDomainID() error {
	url := "https://api.cloudflare.com/client/v4/zones/" + cfg.zoneID + "/dns_records"
	req, err := http.NewRequest("GET", url, nil)
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
	if cfg.records == nil {
		cfg.records = make(map[string]string)
	}

	// Find and store IDs for all configured record names
	for _, record := range dnsInfo.Result {
		if _, exists := cfg.records[record.Name]; exists {
			cfg.records[record.Name] = record.ID
		}
	}

	// Verify we found all records
	for name, id := range cfg.records {
		if id == "" {
			return fmt.Errorf("could not find DNS record for %s", name)
		}
	}

	return nil
}

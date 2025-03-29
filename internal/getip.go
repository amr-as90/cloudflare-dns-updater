package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const url = "https://api.ipify.org?format=json"

type IPResponse struct {
	IP string `json:"ip"`
}

// GetIP fetches the current IP address
func GetIP() (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("ERROR: Could not make query to IPIFY. %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ipResponse IPResponse
	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		return "", err
	}
	return ipResponse.IP, nil
}

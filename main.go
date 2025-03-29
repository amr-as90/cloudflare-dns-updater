package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)


var CurrentIP string

func init() {
	ip, err := GetIP()
	if err != nil {
		fmt.Printf("Could not get new IP: %s", err)
		return
	}

	CurrentIP = ip
}

func main() {
	if err := godotenv.Load(".config"); err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println("Loaded .env file")
	}

	zoneID := os.Getenv("ZONE_ID")
	if zoneID == "" {
		log.Fatal("Zone ID was not set")
	}

	recordID := os.Getenv("DNS_RECORD_ID")
	if recordID == "" {
		log.Fatal("No DNS records set")
	}

	email := os.Getenv("EMAIL")
	if email == "" {
		log.Fatal("No e-mail set")
	}

	apiKey := os.Getenv("ACCESS_TOKEN")
	if apiKey == "" {
		log.Fatal("No API key set")
	}

	cfg := cfConfig{
		zoneID:      zoneID,
		dnsRecordID: recordID,
		eMail:       email,
		apiKey:      apiKey,
	}

	fmt.Printf("Current IP is: %s\n", CurrentIP)
	if CurrentIP == "" {
		fmt.Println("Current IP is empty")
		os.Exit(1)
	}

	cfg.cloudFlareUpdate(CurrentIP)

}

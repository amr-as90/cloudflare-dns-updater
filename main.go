package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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

	recordIDs := os.Getenv("DNS_RECORD_ID")
	if recordIDs == "" {
		log.Fatal("No DNS records set")
	}

	recordIDsSplit := strings.Split(recordIDs, ",")
	for i := range recordIDsSplit {
		recordIDsSplit[i] = strings.TrimSpace(recordIDsSplit[i])
	}

	recordNames := os.Getenv("RECORD_NAMES")
	if recordNames == "" {
		log.Fatal("No record names set")
	}

	recordNamesSplit := strings.Split(recordNames, ",")

	for i := range recordNamesSplit {
		recordNamesSplit[i] = strings.TrimSpace(recordNamesSplit[i])
	}

	if len(recordNamesSplit) != len(recordIDsSplit) {
		log.Fatal("Number of record ID entries must match record name entries")
	}

	email := os.Getenv("EMAIL")
	if email == "" {
		log.Fatal("No e-mail set")
	}

	apiKey := os.Getenv("ACCESS_TOKEN")
	if apiKey == "" {
		log.Fatal("No API key set")
	}

	updateIntervalString := os.Getenv("UPDATE_INTERVAL")
	updateInterval := 60 // Default to 60 seconds if not specified
	if updateIntervalString != "" {
		interval, err := strconv.Atoi(updateIntervalString)
		if err != nil {
			log.Fatalf("Error while getting update interval: %s", err)
		}
		updateInterval = interval
	}

	pushoverAppToken := os.Getenv("PUSHOVER_APP_TOKEN")
	pushoverUserKey := os.Getenv("PUSHOVER_USER_KEY")

	cfg := cfConfig{
		zoneID:           zoneID,
		dnsRecordID:      recordIDsSplit,
		dnsRecordNames:   recordNamesSplit,
		eMail:            email,
		apiKey:           apiKey,
		pushoverAppToken: pushoverAppToken,
		pushoverUserKey:  pushoverUserKey,
	}

	fmt.Printf("Current IP is: %s\n", CurrentIP)
	if CurrentIP == "" {
		fmt.Println("Current IP is empty")
		os.Exit(1)
	}

	// Initial update
	err := cfg.cloudFlareUpdate(CurrentIP)
	if err != nil {
		log.Printf("Error updating DNS records: %s", err)
	}

	ticker := time.NewTicker(time.Duration(updateInterval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("IP check scheduled every %d seconds\n", updateInterval)

	for {
		select {
		case <-ticker.C:
			newIP, err := GetIP()
			if err != nil {
				log.Printf("Error getting IP: %s", err)
				continue
			}

			if newIP != CurrentIP {
				fmt.Printf("IP changed from %s to %s\n", CurrentIP, newIP)
				err = cfg.cloudFlareUpdate(newIP)
				if err != nil {
					log.Printf("Error updating DNS records: %s", err)
				} else {
					CurrentIP = newIP

				}
			}
		}
	}
}

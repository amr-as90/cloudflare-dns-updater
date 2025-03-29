package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amr-as90/cloudflare-dns-updater/internal"

	"github.com/joho/godotenv"
)

var CurrentIP string

func init() {
	ip, err := internal.GetIP()
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

	recordNames := os.Getenv("RECORD_NAMES")
	if recordNames == "" {
		log.Fatal("No record names set")
	}

	email := os.Getenv("EMAIL")
	if email == "" {
		log.Fatal("Email was not set")
	}

	apiKey := os.Getenv("ACCESS_TOKEN")
	if apiKey == "" {
		log.Fatal("API key was not set")
	}

	updateInterval := 60 // Default to 60 seconds if not set

	updateIntervalStr := os.Getenv("UPDATE_INTERVAL")
	if updateIntervalStr != "" {
		interval, err := strconv.Atoi(updateIntervalStr)
		if err != nil {
			log.Fatalf("Invalid update interval: %s", err)
		}
		updateInterval = interval
	}

	pushoverAppToken := os.Getenv("PUSHOVER_APP_TOKEN")
	pushoverUserKey := os.Getenv("PUSHOVER_USER_KEY")

	recordNamesSplit := strings.Split(recordNames, ",")
	records := make(map[string]string)
	for _, name := range recordNamesSplit {
		records[strings.TrimSpace(name)] = "" // Initialize with empty IDs
	}

	cfg := internal.CFConfig{
		ZoneID:           zoneID,
		Records:          records,
		EMail:            email,
		ApiKey:           apiKey,
		PushoverAppToken: pushoverAppToken,
		PushoverUserKey:  pushoverUserKey,
	}

	fmt.Printf("Current IP is: %s\n", CurrentIP)
	if CurrentIP == "" {
		fmt.Println("Current IP is empty")
		os.Exit(1)
	}

	err := cfg.GetDomainID()
	if err != nil {
		log.Fatalf("Error while getting domain ID's: %s", err)
	}

	// Initial update
	err = cfg.CloudFlareUpdate(CurrentIP)
	if err != nil {
		log.Printf("Error updating DNS records: %s", err)
	}

	ticker := time.NewTicker(time.Duration(updateInterval) * time.Second)
	defer ticker.Stop()

	fmt.Printf("Checking IP every %d seconds\n", updateInterval)

	for {
		select {
		case <-ticker.C:
			newIP, err := internal.GetIP()
			if err != nil {
				log.Printf("Error getting IP: %s", err)
				continue
			}

			if newIP != CurrentIP {
				fmt.Printf("IP changed from %s to %s\n", CurrentIP, newIP)
				err = cfg.CloudFlareUpdate(newIP)
				if err != nil {
					log.Printf("Error updating DNS records: %s", err)
				} else {
					CurrentIP = newIP

				}
			}
		}
	}
}

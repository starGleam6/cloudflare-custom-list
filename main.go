package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"time"
)

type Config struct {
	APIKey          string   `yaml:"api_key"`
	APIEmail        string   `yaml:"api_email"`
	AccountID       string   `yaml:"account_id"`
	ListID          string   `yaml:"list_id"`
	DomainNames     []string `yaml:"domain_names"`
	FixedIPs        []string `yaml:"fixed_ips"`
	IntervalMinutes int      `yaml:"interval_minutes"`
	ReplaceList     bool     `yaml:"replace_list"`
}

type IPRecord struct {
	Domain string
	IP     string
}

func getRecords(domainNames []string) ([]IPRecord, error) {
	var records []IPRecord
	for _, domain := range domainNames {
		ipRecords, err := net.LookupIP(domain)
		if err != nil {
			return nil, err
		}

		for _, ip := range ipRecords {
			if ip.To4() != nil { // 检查是否为IPv4地址
				records = append(records, IPRecord{Domain: domain, IP: ip.String()})
			}
		}
	}
	return records, nil
}

func clearIPList(api *cloudflare.API, accountID string, listID string) error {
	items, err := api.ListIPListItems(context.Background(), accountID, listID)
	if err != nil {
		return err
	}

	// If items is empty, nothing to delete
	if len(items) == 0 {
		return nil
	}

	var deleteItems []cloudflare.IPListItemDeleteItemRequest
	for _, item := range items {
		deleteItems = append(deleteItems, cloudflare.IPListItemDeleteItemRequest{ID: item.ID})
	}
	deleteRequest := cloudflare.IPListItemDeleteRequest{Items: deleteItems}

	_, err = api.DeleteIPListItems(context.Background(), accountID, listID, deleteRequest)
	if err != nil {
		return err
	}

	return nil
}

func updateIPList(api *cloudflare.API, accountID string, listID string, records []IPRecord) error {
	for _, record := range records {
		_, err := api.CreateIPListItem(context.Background(), accountID, listID, record.IP, record.Domain)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Updated IP list with %d IPs\n----------\n", len(records))
	return nil
}

func listExists(api *cloudflare.API, accountID string, listID string) (bool, error) {
	lists, err := api.ListIPLists(context.Background(), accountID)
	if err != nil {
		return false, err
	}

	for _, list := range lists {
		if list.ID == listID {
			return true, nil
		}
	}

	return false, nil
}

func performUpdate(api *cloudflare.API, config *Config) {
	fmt.Printf("Current time: %s Performing update...\n----------\n", time.Now().Format("2006-01-02 15:04:05"))

	records, err := getRecords(config.DomainNames)
	if err != nil {
		fmt.Printf("Error looking up IP addresses: %s\n", err)
		return
	}

	for _, ip := range config.FixedIPs {
		record := IPRecord{
			Domain: "",
			IP:     ip,
		}
		records = append(records, record)
	}

	fmt.Printf("Records: %v\n", records)

	exists, err := listExists(api, config.AccountID, config.ListID)
	if err != nil {
		fmt.Printf("Error checking if list exists: %s\n", err)
		return
	}

	if exists {
		fmt.Printf("CloudFlare IP list exists\n----------\n")
		if config.ReplaceList {
			err = clearIPList(api, config.AccountID, config.ListID)
			if err != nil {
				fmt.Printf("Error clearing IP list: %s\n", err)
				return
			}
			fmt.Printf("Cleared IP list\n----------\n")
		}

		// If the list already exists, we update it
		err = updateIPList(api, config.AccountID, config.ListID, records)
		if err != nil {
			fmt.Printf("Error updating IP list: %s\n", err)
		}
	} else {
		// If the list does not exist, we create it
		fmt.Printf("IP list does not exist\n----------\n")
		_, err = api.CreateIPList(context.Background(), config.AccountID, config.ListID, "auto-updated list", "ip")
		if err != nil {
			fmt.Printf("Error creating IP list: %s\n", err)
			return
		}

		err = updateIPList(api, config.AccountID, config.ListID, records)
		if err != nil {
			fmt.Printf("Error updating newly created IP list: %s\n", err)
		}
	}
	fmt.Printf("Update complete\n----------\n")
}

func main() {
	// Open log file
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer logFile.Close()

	// Redirect standard output to log file
	oldStdout := os.Stdout
	os.Stdout = logFile
	defer func() { os.Stdout = oldStdout }()

	// Read and parse configuration file
	fmt.Printf("----------\nReading config file...\n----------\n")

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		return
	}
	fmt.Printf("data: %s\n", data)
	fmt.Printf("config: %s\n", config)

	fmt.Printf("Readed config file\n----------\n")

	api, err := cloudflare.New(config.APIKey, config.APIEmail)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Perform an immediate update
	performUpdate(api, &config)

	ticker := time.NewTicker(time.Duration(config.IntervalMinutes) * time.Minute) // 定期更新的时间间隔

	// Perform periodic updates
	for range ticker.C {
		performUpdate(api, &config)
	}
}

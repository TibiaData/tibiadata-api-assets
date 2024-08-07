package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// sleepFlag should be to sleep between requests
	sleepFlag bool

	TibiaComHost     = "www.tibia.com"
	TibiaComProtocol = "https"
)

const (
	TibiaDataAPIhost = "api.tibiadata.com"
)

func init() {
	flag.BoolVar(&sleepFlag, "sleep", false, "Set to sleep between requests")

	// Parse the flags
	flag.Parse()
}

func main() {
	// logging start of TibiaData
	log.Printf("[info] TibiaData assets generator starting..")

	// Setting up resty client
	client := resty.New()

	// Set client timeout  and retry
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	// Set headers for all requests
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "TibiaData assets-generator",
	})

	// Enabling Content length value for all request
	client.SetContentLength(true)

	// Disable redirection of client (so we skip parsing maintenance page)
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	// overriding host with env
	if isEnvExist("TIBIADATA_PROXY") {
		TibiaComHost = getEnv("TIBIADATA_PROXY", "www.tibia.com")

		// overriding protocol with env
		if isEnvExist("TIBIADATA_PROXY_PROTOCOL") {
			TibiaComProtocol = getEnv("TIBIADATA_PROXY_PROTOCOL", "https")
		}
	}

	var builder Builder

	err := builder.housesWorker(client)
	if err != nil {
		log.Fatalf("[error] Issue with housesWorker. Error: %s", err)
	}

	if sleepFlag {
		time.Sleep(time.Second / 2)
	}

	err = builder.creaturesWorker(client)
	if err != nil {
		log.Fatalf("[error] Issue with creaturesWorker. Error: %s", err)
	}

	if sleepFlag {
		time.Sleep(time.Second / 2)
	}

	err = builder.spellsWorker(client)
	if err != nil {
		log.Fatalf("[error] Issue with fansitesWorker. Error: %s", err)
	}

	log.Println("[info] Validation of builder lists to prevent empty set of strings.")
	validateBuilderLists(builder)
	log.Println("[info] Validation of builder lists passed.")

	log.Println("[info] Generating output file: output.json")
	outputFile, err := json.Marshal(builder)
	if err != nil {
		log.Fatalf("[error] Issue with marshaling output file. Error: %s", err)
	}

	err = os.WriteFile("output.json", outputFile, 0644)
	if err != nil {
		log.Fatalf("[error] Issue writing the min file. Error: %s", err)
	}

	log.Println("[info] TibiaData assets generator finished.")
}

// isEnvExist func - check if environment var is set
func isEnvExist(key string) (ok bool) {
	_, ok = os.LookupEnv(key)
	return
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}

	return defaultVal
}

// validateBuilderLists func - check if any set of strings is empty
func validateBuilderLists(builder Builder) {
	v := reflect.ValueOf(builder)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Slice && field.Len() == 0 {
			fieldName := v.Type().Field(i).Name
			log.Fatalf("[error] Validation of builder lists failed. Empty list: %s.\n", fieldName)
		}
	}
}

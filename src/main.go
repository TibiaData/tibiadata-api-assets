package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"
)

const (
	TibiaDataAPIhost = "dev.tibiadata.com"
)

type AssetsHouse struct {
	HouseID   int    `json:"house_id"`
	Town      string `json:"town"`
	HouseType string `json:"type"`
}

type AssetsHouses struct {
	Worlds []string      `json:"worlds"`
	Towns  []string      `json:"towns"`
	Houses []AssetsHouse `json:"houses"`
}

type SourceHousesOverview struct {
	Houses struct {
		HouseList []struct {
			HouseID int `json:"house_id"`
		} `json:"house_list"`
		GuildhallList []struct {
			HouseID int `json:"house_id"`
		} `json:"guildhall_list"`
	} `json:"houses"`
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

	// defining values for request
	var (
		res *resty.Response
		err error

		AssetsHouses AssetsHouses

		TibiaComhost = "www.tibia.com"
	)

	// overriding host with env
	if isEnvExist("TIBIADATA_PROXY") {
		TibiaComhost = "https://" + getEnv("TIBIADATA_PROXY", "www.tibia.com")
	}

	res, err = client.R().Get(TibiaComhost + "/community/?subtopic=houses")

	switch res.StatusCode() {
	case http.StatusOK:
		log.Println("[info] Retrieving data successfully from tibia.com.")
	default:
		log.Fatalf("[error] Issue when collecting data from tibia.com. Error: %s", err)
	}

	// Convert body to io.Reader
	resIo := bytes.NewReader(res.Body())
	// wrap reader in a converting reader from ISO 8859-1 to UTF-8
	resIo2 := norm.NFKC.Reader(charmap.ISO8859_1.NewDecoder().Reader(resIo))

	// Load the HTML document
	doc, _ := goquery.NewDocumentFromReader(resIo2)
	if err != nil {
		log.Fatalf("[error] Issue with goquery reading document. Error: %s", err)
	}

	// Find of this to get div with class BoxContent
	doc.Find(".TableContentContainer .TableContent tbody tr").First().Next().Children().Each(func(index int, s *goquery.Selection) {
		// generate list of worlds that have houses/guildhalls
		s.Find("select").Children().NextAll().Each(func(i int, selection *goquery.Selection) {
			// collect the world
			AssetsHouses.Worlds = append(AssetsHouses.Worlds, selection.Text())
		})

		// generate list of towns that have houses/guildhalls
		s.Find("input[name=town]").Each(func(i int, selection *goquery.Selection) {
			// collect the town
			AssetsHouses.Towns = append(AssetsHouses.Towns, selection.AttrOr("value", ""))
		})

	})

	for _, town := range AssetsHouses.Towns {
		log.Printf("[info] Retrieving data about houses and guildhalls in %s.", town)

		// sleep for 500 ms for ratelimit on dev.tibiadata.com
		time.Sleep(500 * time.Millisecond)

		ApiUrl := "https://" + TibiaDataAPIhost + "/v3/houses/" + AssetsHouses.Worlds[0] + "/" + url.QueryEscape(town)
		res, err = client.R().Get(ApiUrl)

		switch res.StatusCode() {
		case http.StatusOK:
			// Get byte slice from string.
			bytes := []byte(res.Body())

			var cont SourceHousesOverview
			err := json.Unmarshal(bytes, &cont)
			if err != nil {
				log.Fatalf("[error] Issue when unmarshaling data. Town is %s. Err: %s", town, err)
			}

			for _, value := range cont.Houses.HouseList {
				AssetsHouses.Houses = append(AssetsHouses.Houses, AssetsHouse{
					HouseID:   value.HouseID,
					Town:      town,
					HouseType: "house",
				})
			}
			for _, value := range cont.Houses.GuildhallList {
				AssetsHouses.Houses = append(AssetsHouses.Houses, AssetsHouse{
					HouseID:   value.HouseID,
					Town:      town,
					HouseType: "guildhall",
				})
			}

		default:
			log.Fatalf("[error] Issue when collecting data from %s. Error: %s", TibiaDataAPIhost, err)
		}

	}

	log.Println("[info] Generating output file: docs/data.json")
	file, _ := json.MarshalIndent(AssetsHouses, "", " ")
	_ = ioutil.WriteFile("docs/data.json", file, 0644)

	log.Println("[info] TibiaData assets generator finished.")
}

// isEnvExist func - check if environment var is set
func isEnvExist(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}

	return false
}

// getEnv func - read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

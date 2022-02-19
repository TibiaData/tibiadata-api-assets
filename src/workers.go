package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/unicode/norm"
)

func getDoc(client *resty.Client, endpoint string) (*goquery.Document, error) {
	res, err := client.R().Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Issue getting doc for %s endpoint. Error: %s", endpoint, err)
	}

	switch res.StatusCode() {
	case http.StatusOK:
		log.Printf("[info] Retrieving data successfully from tibia.com. Endpoint: %s\n", endpoint)
	default:
		return nil, fmt.Errorf("Issue when collecting data from tibia.com. StatusCode: %d, Endpoint: %s", res.StatusCode(), endpoint)
	}

	// Convert body to io.Reader
	resIo := bytes.NewReader(res.Body())
	// wrap reader in a converting reader from ISO 8859-1 to UTF-8
	resIo2 := norm.NFKC.Reader(charmap.ISO8859_1.NewDecoder().Reader(resIo))

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resIo2)
	if err != nil {
		return nil, fmt.Errorf("Issue with goquery reading document. Error: %s, Endpoint: %s", err, endpoint)
	}

	return doc, nil
}

func (b *Builder) housesWorker(client *resty.Client) error {
	doc, err := getDoc(client, "https://"+TibiaComHost+"/community/?subtopic=houses")
	if err != nil {
		return err
	}

	// Find of this to get div with class BoxContent
	doc.Find(".TableContentContainer .TableContent tbody tr").First().Next().Children().Each(func(index int, s *goquery.Selection) {
		// generate list of worlds that have houses/guildhalls
		s.Find("select").Children().NextAll().Each(func(i int, selection *goquery.Selection) {
			// collect the world
			b.Worlds = append(b.Worlds, selection.Text())
		})

		// generate list of towns that have houses/guildhalls
		s.Find("input[name=town]").Each(func(i int, selection *goquery.Selection) {
			// collect the town
			b.Towns = append(b.Towns, selection.AttrOr("value", ""))
		})

	})

	for _, town := range b.Towns {
		log.Printf("[info] Retrieving data about houses and guildhalls in %s.", town)

		ApiUrl := "https://" + TibiaDataAPIhost + "/v3/houses/" + b.Worlds[0] + "/" + url.QueryEscape(town)
		res, err := client.R().Get(ApiUrl)
		if err != nil {
			return fmt.Errorf("Issue getting %s endpoint. Error: %s", ApiUrl, err)
		}

		switch res.StatusCode() {
		case http.StatusOK:
			// Get byte slice from string.
			bytes := []byte(res.Body())

			var cont SourceHousesOverview
			err := json.Unmarshal(bytes, &cont)
			if err != nil {
				return fmt.Errorf("Issue when unmarshaling data. Town is %s. Err: %s", town, err)
			}

			for _, value := range cont.Houses.HouseList {
				b.Houses = append(b.Houses, AssetsHouse{
					HouseID:   value.HouseID,
					Town:      town,
					HouseType: "house",
				})
			}

			for _, value := range cont.Houses.GuildhallList {
				b.Houses = append(b.Houses, AssetsHouse{
					HouseID:   value.HouseID,
					Town:      town,
					HouseType: "guildhall",
				})
			}
		default:
			return fmt.Errorf("Issue when collecting data from %s. StatusCode: %d", TibiaDataAPIhost, res.StatusCode())
		}

		if sleepFlag {
			time.Sleep(time.Second / 2)
		}
	}

	return nil
}

func (b *Builder) creaturesWorker(client *resty.Client) error {
	doc, err := getDoc(client, "https://"+TibiaComHost+"/library/?subtopic=creatures")
	if err != nil {
		return err
	}

	const raceEndpoint = "https://www.tibia.com/library/?subtopic=creatures&race="

	var safe []string

	creatures := doc.Find(".BoxContent div[style] div[style]").First()
	creatures.Find("div").Each(func(index int, s *goquery.Selection) {
		url, exists := s.Find("a").Attr("href")
		if !exists {
			return
		}

		endpoint := strings.TrimPrefix(url, raceEndpoint)
		safe = append(safe, endpoint)
		pluralName := s.Find("div").First().Text()
		fields := strings.Fields(pluralName)
		length := len(fields)

		if name, ok := specialCreaturesCases[endpoint]; ok {
			b.Creatures = append(b.Creatures, Creature{
				Endpoint:   endpoint,
				PluralName: pluralName,
				Name:       name,
			})
		} else if length == 1 {
			b.Creatures = append(b.Creatures, Creature{
				Endpoint:   endpoint,
				PluralName: pluralName,
				Name:       strings.Title(endpoint),
			})
		} else {
			var rawNameBuilder string
			var failed int

			for i, f := range fields {
				currentField := strings.ToLower(f)

				index := strings.Index(endpoint, currentField)
				if index == -1 {
					failed = i
					continue
				}

				rawNameBuilder += currentField
			}

			missingWord := strings.Title(strings.ReplaceAll(endpoint, rawNameBuilder, ""))

			var creature string
			for i, f := range fields {
				if i == failed {
					creature += missingWord
				} else {
					creature += f
				}

				if i != length-1 {
					creature += " "
				}
			}

			b.Creatures = append(b.Creatures, Creature{
				Endpoint:   endpoint,
				PluralName: pluralName,
				Name:       creature,
			})
		}
	})

	for i, s := range safe {
		str := SpaceMap(b.Creatures[i].Name)
		_, isSpecial := specialCreaturesCases[s]

		if !isSpecial && !strings.EqualFold(s, str) {
			return fmt.Errorf("[error] Wrong creature name passed. Expected: %s, got: %s", s, str)
		}
	}

	return nil
}

func (b *Builder) spellsWorker(client *resty.Client) error {
	doc, err := getDoc(client, "https://"+TibiaComHost+"/library/?subtopic=spells")
	if err != nil {
		return err
	}

	doc.Find("table.TableContent ~ table tr").First().NextAll().Each(func(index int, s *goquery.Selection) {
		if index == 0 {
			return
		}

		s.Find("td").EachWithBreak(func(index int, inner *goquery.Selection) bool {
			if index == 0 {
				rawText := inner.Text()
				spellName := rawText[0:strings.Index(rawText, " (")]
				formula := rawText[strings.Index(rawText, " (")+2 : strings.Index(rawText, ")")]

				var endpoint string
				if specialCase, isSpecial := specialSpellsCases[spellName]; isSpecial {
					endpoint = specialCase
				} else {
					endpoint = SanitizeSpellEndpoint(spellName)
				}

				b.Spells = append(b.Spells, Spell{
					Name:     spellName,
					Formula:  formula,
					Endpoint: endpoint,
				})

				return false
			}

			return true
		})
	})

	return nil
}

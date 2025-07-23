package main

import (
	"strings"
	"unicode"
)

var (
	// Special creature name cases
	specialCreaturesCases = map[string]string{
		"apparitionofadruid":     "Druid's Apparition",
		"asura":                  "Dawnfire Asura",
		"asuranight":             "Midnight Asura",
		"batface":                "Gloom Maws",
		"caribbeanbat":           "Exotic Bat",
		"caribbeancavespider":    "Exotic Cave Spider",
		"carnisylvandark":        "Dark Carnisylvan",
		"carnisylvanhulking":     "Hulking Carnisylvan",
		"carnisylvanpoisonous":   "Poisonous Carnisylvan",
		"charlatan":              "Corym Charlatan",
		"crustaceagigantica":     "Crustacea Gigantica",
		"crystalgolem":           "Enraged Crystal Golem",
		"cultacolyte":            "Acolyte Of The Cult",
		"cultadept":              "Adept Of The Cult",
		"cultnovice":             "Novice Of The Cult",
		"cultpriest":             "Enlightened Of The Cult",
		"degeneratedshaper":      "Broken Shaper",
		"earthelementalmassive":  "Massive Earth Elemental",
		"energyelementalmassive": "Massive Energy Elemental",
		"frogazure":              "Azure Frog",
		"ghostlycrawler":         "Branchy Crawler",
		"girtablilu":             "Venerable Girtablilu",
		"hellfireelemental":      "Massive Fire Elemental",
		"knightsapparition":      "Knight's Apparition",
		"lionmonk":               "Monk Of The Order",
		"lostsoulhard":           "Freakish Lost Soul",
		"lostsoulmedium":         "Mean Lost Soul",
		"lostsoulweak":           "Flimsy Lost Soul",
		"manyfaces":              "Many Faces",
		"monksapparition":        "Monk's Apparition",
		"moohtahwarrior":         "Mooh'tah Warrior",
		"norcferatudworc":        "Dworc Shadowstalkers",
		"norcferatuorclops":      "Orclops Bloodbreakers",
		"paladinsapparition":     "Paladin's Apparition",
		"ragingbrainsquid":       "Rage Squid",
		"skirmisher":             "Corym Skirmisher",
		"sorcerersapparition":    "Sorcerer's Apparition",
		"soulbrokenharbinger":    "Soul-broken Harbinger",
		"twoheadedturtle":        "Two-headed Turtles",
		"vanguard":               "Corym Vanguard",
		"viscountmanbat":         "Vicious Manbat",
		"waterelementalmassive":  "Massive Water Elemental",
		"whitedeer":              "White Deer",
		"wraith":                 "Betrayed Wraith",
	}

	specialSpellsCases = map[string]string{
		"Apprentice's Strike": "apprenticestrike",
	}
)

func SanitizeSpellEndpoint(spell string) string {
	s := strings.ToLower(spell)
	s = strings.ReplaceAll(s, "'", "")
	s = SpaceMap(s)
	return s
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

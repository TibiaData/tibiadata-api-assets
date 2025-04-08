package main

import (
	"strings"
	"unicode"
)

var (
	// Special creature name cases
	specialCreaturesCases = map[string]string{
		"cultacolyte":            "Acolyte Of The Cult",
		"cultadept":              "Adept Of The Cult",
		"frogazure":              "Azure Frog",
		"wraith":                 "Betrayed Wraith",
		"ghostlycrawler":         "Branchy Crawler",
		"degeneratedshaper":      "Broken Shaper",
		"charlatan":              "Corym Charlatan",
		"skirmisher":             "Corym Skirmisher",
		"vanguard":               "Corym Vanguard",
		"crustaceagigantica":     "Crustacea Gigantica",
		"carnisylvandark":        "Dark Carnisylvan",
		"asura":                  "Dawnfire Asura",
		"apparitionofadruid":     "Druid's Apparition",
		"cultpriest":             "Enlightened Of The Cult",
		"crystalgolem":           "Enraged Crystal Golem",
		"caribbeanbat":           "Exotic Bat",
		"caribbeancavespider":    "Exotic Cave Spider",
		"lostsoulweak":           "Flimsy Lost Soul",
		"lostsoulhard":           "Freakish Lost Soul",
		"carnisylvanhulking":     "Hulking Carnisylvan",
		"knightsapparition":      "Knight's Apparition",
		"manyfaces":              "Many Faces",
		"earthelementalmassive":  "Massive Earth Elemental",
		"energyelementalmassive": "Massive Energy Elemental",
		"hellfireelemental":      "Massive Fire Elemental",
		"waterelementalmassive":  "Massive Water Elemental",
		"lostsoulmedium":         "Mean Lost Soul",
		"asuranight":             "Midnight Asura",
		"lionmonk":               "Monk Of The Order",
		"monksapparition":        "Monk's Apparition",
		"moohtahwarrior":         "Mooh'tah Warrior",
		"cultnovice":             "Novice Of The Cult",
		"paladinsapparition":     "Paladin's Apparition",
		"carnisylvanpoisonous":   "Poisonous Carnisylvan",
		"ragingbrainsquid":       "Rage Squid",
		"sorcerersapparition":    "Sorcerer's Apparition",
		"soulbrokenharbinger":    "Soul-broken Harbinger",
		"twoheadedturtle":        "Two-headed Turtles",
		"girtablilu":             "Venerable Girtablilu",
		"viscountmanbat":         "Vicious Manbat",
		"whitedeer":              "White Deer",
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

package main

type AssetsHouse struct {
	Name      string `json:"name"`
	HouseID   int    `json:"house_id"`
	Town      string `json:"town"`
	HouseType string `json:"type"`
}

type SourceHousesOverview struct {
	Houses Houses `json:"houses"`
}

type Houses struct {
	HouseList     []House `json:"house_list"`
	GuildhallList []House `json:"guildhall_list"`
}

type House struct {
	Name    string `json:"name"`
	HouseID int    `json:"house_id"`
}

type Creature struct {
	Endpoint   string `json:"endpoint"`
	PluralName string `json:"plural_name"`
	Name       string `json:"name"`
}

type Spell struct {
	Name     string `json:"name"`
	Formula  string `json:"formula"`
	Endpoint string `json:"endpoint"`
}

type Builder struct {
	Worlds    []string      `json:"worlds"`
	Towns     []string      `json:"towns"`
	Houses    []AssetsHouse `json:"houses"`
	Creatures []Creature    `json:"creatures"`
	Spells    []Spell       `json:"spells"`
}

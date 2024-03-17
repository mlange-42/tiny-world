package save

import (
	"image"
	"time"
)

const mapDescriptionDelimiter = "----"

type LoadType uint8

const (
	LoadTypeNone LoadType = iota
	LoadTypeGame
	LoadTypeMap
)

type MapLocation struct {
	Name       string
	IsEmbedded bool
}

type SaveGame struct {
	Name string
	Time time.Time
}

type saveGameInfo struct {
	Resources saveGameResources
}

type saveGameResources struct {
	SaveTime saveTime `json:"res.SaveTime"`
}

type saveTime struct {
	Time time.Time
}

type MapInfo struct {
	Achievements []string
	Description  string
}

type mapJs struct {
	Terrains              map[string]int `json:"terrains"`
	Map                   []string       `json:"map"`
	Achievements          []string       `json:"achievements"`
	Description           []string       `json:"description"`
	Center                image.Point    `json:"center"`
	InitialRandomTerrains int            `json:"initial_terrains"`
}

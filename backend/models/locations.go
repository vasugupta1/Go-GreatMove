package models

type TypeAheadLocation struct {
	DisplayName          string `json:"displayName"`
	LocationIdentifier   string `json:"locationIdentifier"`
	NormalisedSearchTerm string `json:"normalisedSearchTerm"`
}

type RightMoveLocationResponse struct {
	Key                string              `json:"key"`
	Term               string              `json:"term"`
	TypeAheadLocations []TypeAheadLocation `json:"typeAheadLocations"`
}

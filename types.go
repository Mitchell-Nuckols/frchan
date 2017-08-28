package main

// TBATeam : GET /team/{team_key}
type TBATeam struct {
	Key              string   `json:"key"`
	TeamNumber       int      `json:"team_number"`
	Nickname         string   `json:"nickname"`
	Name             string   `json:"name"`
	City             string   `json:"city"`
	SateProv         string   `json:"state_prov"`
	Country          string   `json:"country"`
	Address          string   `json:"address,omitempty"`
	PostalCode       string   `json:"postal_code"`
	GMapsPlaceID     string   `json:"gmaps_place_id,omitempty"`
	GMapsPlaceURL    string   `json:"gmaps_url,omitempty"`
	Lat              int      `json:"lat"`
	Lng              int      `json:"lng"`
	LocationName     string   `json:"location_name,omitempty"`
	Website          string   `json:"website,omitempty"`
	RookieYear       int      `json:"rookie_year"`
	Motto            string   `json:"motto,omitempty"`
	HomeChampionship struct{} `json:"home_championship,omitempty"`
}

// TBATeamAwards : GET /team/{team_key}/awards
type TBATeamAwards []struct {
	Name          string `json:"name"`
	AwardType     int    `json:"award_type"`
	EventKey      string `json:"event_key"`
	RecipientList []struct {
		TeamKey string `json:"team_key"`
		Awardee string `json:"awardee"`
	} `json:"recipient_list"`
	Year int `json:"year"`
}

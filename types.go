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

// TBAEventRankings : GET /event/{event_key}/rankings
type TBAEventRankings struct {
	Rankings []struct {
		Dq            int `json:"dq"`
		MatchesPlayed int `json:"matches_played"`
		QualAverage   int `json:"qual_average"`
		Rank          int `json:"rank"`
		Record        struct {
			Losses int `json:"losses"`
			Wins   int `json:"wins"`
			Ties   int `json:"ties"`
		} `json:"record"`
		SortOrders []int  `json:"sort_orders"`
		TeamKey    string `json:"team_key"`
	} `json:"rankings"`
	SortOrderInfo []struct {
		Name      string `json:"name"`
		Precision int    `json:"precision"`
	} `json:"sort_order_info"`
}

// TBATeamEventsSimple : GET /team/{team_key}/events/simple
type TBATeamEventsSimple []struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	EventCode string `json:"event_code"`
	EventType int    `json:"event_type"`
	District  struct {
		Abbreviation string `json:"abbreviation"`
		DisplayName  string `json:"display_name"`
		Key          string `json:"key"`
		Year         int    `json:"year"`
	} `json:"district"`
	City      string `json:"city"`
	StateProv string `json:"state_prov"`
	Country   string `json:"country"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Year      int    `json:"year"`
}

// TBAEventTeamKeys : GET /event/{event_key}/teams/keys
type TBAEventTeamKeys []string

// TBATeamEventStatus : GET /team/{team_key}/event/{event_key}/status
type TBATeamEventStatus struct {
	Qual struct {
		NumTeams int `json:"num_teams"`
		Ranking  struct {
			Dq            int `json:"dq"`
			MatchesPlayed int `json:"matches_played"`
			QualAverage   int `json:"qual_average"`
			Rank          int `json:"rank"`
			Record        struct {
				Losses int `json:"losses"`
				Wins   int `json:"wins"`
				Ties   int `json:"ties"`
			} `json:"record"`
			SortOrders []float64 `json:"sort_orders"`
			TeamKey    string    `json:"team_key"`
		} `json:"ranking"`
		SortOrderInfo []struct {
			Name      string `json:"name"`
			Precision int    `json:"precision"`
		} `json:"sort_order_info"`
		Status string `json:"status"`
	} `json:"qual"`
	Alliance struct {
		Name   string `json:"name"`
		Number int    `json:"number"`
		Backup struct {
			Out string `json:"out"`
			In  string `json:"in"`
		} `json:"backup"`
		Pick int `json:"pick"`
	} `json:"alliance"`
	Playoff           TBAPlayoff `json:"playoff"`
	AllianceStatusStr string     `json:"alliance_status_str"`
	PlayoffStatusStr  string     `json:"playoff_status_str"`
	OverallStatusStr  string     `json:"overall_status_str"`
}

// TBAPlayoff : Playoff info struct
type TBAPlayoff struct {
	Level              string `json:"level"`
	CurrentLevelRecord struct {
		Losses int `json:"losses"`
		Wins   int `json:"wins"`
		Ties   int `json:"ties"`
	} `json:"current_level_record"`
	Record struct {
		Losses int `json:"losses"`
		Wins   int `json:"wins"`
		Ties   int `json:"ties"`
	} `json:"record"`
	Status         string `json:"status"`
	PlayoffAverage int    `json:"playoff_average"`
}

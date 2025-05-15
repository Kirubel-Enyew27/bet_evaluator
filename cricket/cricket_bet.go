package cricket

// CricketPrematch represents the prematch data structure
type CricketPrematch struct {
	Success int `json:"success"`
	Results []struct {
		EventID string `json:"event_id"`
		Main    struct {
			UpdatedAt string `json:"updated_at"`
			SP        struct {
				ToWinTheMatch struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Odds []Odd  `json:"odds"`
				} `json:"to_win_the_match"`
			} `json:"sp"`
		} `json:"main"`
		Match struct {
			SP struct {
				MostMatchSixes struct {
					Odds []Odd `json:"odds"`
				} `json:"most_match_sixes"`
				MostMatchFours struct {
					Odds []Odd `json:"odds"`
				} `json:"most_match_fours"`
			} `json:"sp"`
		} `json:"match"`
	} `json:"results"`
}

// Odd represents a betting odd
type Odd struct {
	ID       string `json:"id"`
	Odds     string `json:"odds"`
	Name     string `json:"name"`
	Header   string `json:"header"`
	Handicap string `json:"handicap"`
}

// CricketResult represents the match result data
type CricketResult struct {
	Success int `json:"success"`
	Results []struct {
		ID          string `json:"id"`
		SportID     string `json:"sport_id"`
		Time        string `json:"time"`
		TimeStatus  string `json:"time_status"` // 3 = match completed
		League      League `json:"league"`
		Home        Team   `json:"home"`
		Away        Team   `json:"away"`
		SS          string `json:"ss"` // Score string "117-217" (home-away)
		Extra       Extra  `json:"extra"`
		HasLineup   int    `json:"has_lineup"`
		ConfirmedAt string `json:"confirmed_at"`
	} `json:"results"`
}

type League struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	CC   string `json:"cc"`
}

type Team struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	ImageID string `json:"image_id"`
	CC      string `json:"cc"`
}

type Extra struct {
	StadiumData Stadium `json:"stadium_data"`
}

type Stadium struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Capacity     string `json:"capacity"`
	GoogleCoords string `json:"googlecoords"`
}

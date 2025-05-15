package cricket

import (
	"encoding/json"
	"fmt"
	"os"
)

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

func loadPrematchData(filename string) (*CricketPrematch, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var prematch CricketPrematch
	if err := json.Unmarshal(data, &prematch); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if prematch.Success != 1 || len(prematch.Results) == 0 {
		return nil, fmt.Errorf("invalid or empty prematch data")
	}

	return &prematch, nil
}

func loadResultData(filename string) (*CricketResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var result CricketResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if result.Success != 1 || len(result.Results) == 0 {
		return nil, fmt.Errorf("invalid or empty result data")
	}

	return &result, nil
}

func findMatchingResult(prematch *CricketPrematch, result *CricketResult) (*struct {
	ID          string `json:"id"`
	SportID     string `json:"sport_id"`
	Time        string `json:"time"`
	TimeStatus  string `json:"time_status"`
	League      League `json:"league"`
	Home        Team   `json:"home"`
	Away        Team   `json:"away"`
	SS          string `json:"ss"`
	Extra       Extra  `json:"extra"`
	HasLineup   int    `json:"has_lineup"`
	ConfirmedAt string `json:"confirmed_at"`
}, error) {
	prematchEventID := prematch.Results[0].EventID

	for _, r := range result.Results {
		if r.ID == prematchEventID {
			if r.TimeStatus != "3" {
				return nil, fmt.Errorf("match not completed yet, status: %s", r.TimeStatus)
			}
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no matching result found for event ID %s", prematchEventID)
}

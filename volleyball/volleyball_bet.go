package volleyball

import (
	"encoding/json"
	"log"
	"os"
)

// PrematchResult represents the prematch odds data
type PrematchResult struct {
	FI      string `json:"FI"`
	EventID string `json:"event_id"`
	Main    struct {
		SP struct {
			GameLines struct {
				Odds []struct {
					ID       string `json:"id"`
					Odds     string `json:"odds"`
					Name     string `json:"name"`
					Header   string `json:"header"`
					Handicap string `json:"handicap"`
				} `json:"odds"`
			} `json:"game_lines"`
			CorrectSetScore struct {
				Odds []struct {
					ID     string `json:"id"`
					Odds   string `json:"odds"`
					Name   string `json:"name"`
					Header string `json:"header"`
				} `json:"odds"`
			} `json:"correct_set_score"`
			MatchTotal struct {
				Odds []struct {
					ID       string `json:"id"`
					Odds     string `json:"odds"`
					Name     string `json:"name"`
					Handicap string `json:"handicap"`
				} `json:"odds"`
			} `json:"match_total_odd_even"`
		} `json:"sp"`
	} `json:"main"`
}

// MatchResult represents a single match result
type MatchResult struct {
	EventID string `json:"id"`
	Home    struct {
		Name string `json:"name"`
	} `json:"home"`
	Away struct {
		Name string `json:"name"`
	} `json:"away"`
	SS     string `json:"ss"`
	Scores map[string]struct {
		Home string `json:"home"`
		Away string `json:"away"`
	} `json:"scores"`
}

// ResultData contains multiple match results
type ResultData struct {
	Results []MatchResult `json:"results"`
}

func loadPrematchData(filename string) PrematchResult {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading prematch file: %v", err)
	}

	var data struct {
		Results []PrematchResult `json:"results"`
	}
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatalf("Error parsing prematch JSON: %v", err)
	}

	return data.Results[0]
}

func loadResultData(filename string) ResultData {
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading result file: %v", err)
	}

	var data ResultData
	if err := json.Unmarshal(file, &data); err != nil {
		log.Fatalf("Error parsing result JSON: %v", err)
	}

	return data
}

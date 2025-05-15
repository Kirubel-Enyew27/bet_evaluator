package models

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

// PrematchData contains multiple pre match data
type PrematchData struct {
	Results []PrematchResult `json:"results"`
}

// ResultData contains multiple match results
type ResultData struct {
	Results []MatchResult `json:"results"`
}

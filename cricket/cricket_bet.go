package cricket

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func Evaluate() {
	// Load and parse prematch data
	prematch, err := loadPrematchData("cricket_prematch.json")
	if err != nil {
		log.Fatalf("Error loading prematch data: %v", err)
	}

	// Load and parse result data
	result, err := loadResultData("cricket_result.json")
	if err != nil {
		log.Fatalf("Error loading result data: %v", err)
	}

	// Find matching result for our prematch data
	matchResult, err := findMatchingResult(prematch, result)
	if err != nil {
		log.Fatal(err)
	}

	// Parse scores
	homeScore, awayScore, err := parseScores(matchResult.SS)
	if err != nil {
		log.Fatalf("Error parsing scores: %v", err)
	}

	// Evaluate bets
	evaluateBets(prematch, matchResult, homeScore, awayScore)
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

func parseScores(scoreStr string) (int, int, error) {
	parts := strings.Split(scoreStr, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid score format: %s", scoreStr)
	}

	homeScore, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid home score: %v", err)
	}

	awayScore, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid away score: %v", err)
	}

	return homeScore, awayScore, nil
}

func evaluateBets(prematch *CricketPrematch, result *struct {
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
}, homeScore, awayScore int) {
	fmt.Printf("\n=== Match Evaluation ===\n")
	fmt.Printf("League: %s\n", result.League.Name)
	fmt.Printf("Teams: %s vs %s\n", result.Home.Name, result.Away.Name)
	fmt.Printf("Venue: %s, %s\n", result.Extra.StadiumData.Name, result.Extra.StadiumData.City)
	fmt.Printf("Result: %s (%d - %d)\n", result.SS, homeScore, awayScore)

	fmt.Printf("\n--- Bet Evaluation Results ---\n")

	// Evaluate match winner market
	if len(prematch.Results[0].Main.SP.ToWinTheMatch.Odds) > 0 {
		fmt.Println("\nMatch Winner Market:")
		for _, odd := range prematch.Results[0].Main.SP.ToWinTheMatch.Odds {
			fmt.Printf("- %s @ %s => ", odd.Name, odd.Odds)
			if (odd.Name == "1" && homeScore > awayScore) || (odd.Name == "2" && awayScore > homeScore) {
				fmt.Println("WON")
			} else {
				fmt.Println("LOST")
			}
		}
	}

	// Evaluate most sixes market
	if len(prematch.Results[0].Match.SP.MostMatchSixes.Odds) > 0 {
		fmt.Println("\nMost Match Sixes Market:")
		for _, odd := range prematch.Results[0].Match.SP.MostMatchSixes.Odds {
			fmt.Printf("- %s @ %s => ", odd.Name, odd.Odds)
			if odd.Name == "1" && homeScore > awayScore { // Assuming home team hit more sixes if they scored more runs
				fmt.Println("WON")
			} else if odd.Name == "2" && awayScore > homeScore {
				fmt.Println("WON")
			} else if odd.Name == "Tie" && homeScore == awayScore {
				fmt.Println("WON")
			} else {
				fmt.Println("LOST")
			}
		}
	}

	// Evaluate most fours market
	if len(prematch.Results[0].Match.SP.MostMatchFours.Odds) > 0 {
		fmt.Println("\nMost Match Fours Market:")
		for _, odd := range prematch.Results[0].Match.SP.MostMatchFours.Odds {
			fmt.Printf("- %s @ %s => ", odd.Name, odd.Odds)
			if odd.Name == "1" && homeScore > awayScore {
				fmt.Println("WON")
			} else if odd.Name == "2" && awayScore > homeScore {
				fmt.Println("WON")
			} else if odd.Name == "Tie" && homeScore == awayScore {
				fmt.Println("WON")
			} else {
				fmt.Println("LOST")
			}
		}
	}
}

package volleyball

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

func Evaluate() {
	prematch := loadPrematchData("volleyball/volleyball_prematch.json")
	result := loadResultData("volleyball/volleyball_result.json")

	if len(result.Results) == 0 {
		log.Fatal("No result data found")
	}

	match := result.Results[0]
	fmt.Printf("\nMatch: %s vs %s\n", match.Home.Name, match.Away.Name)
	fmt.Printf("Final Score: %s\n\n", match.SS)

	evaluateWinnerMarket(prematch, match)
	evaluateCorrectScoreMarket(prematch, match)
	evaluateTotalPointsMarket(prematch, match)
	evaluateHandicapMarket(prematch, match)
	evaluateDoubleChanceMarket(prematch, match)

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

func evaluateWinnerMarket(pm PrematchResult, result MatchResult) {
	fmt.Println("--- Match Winner (1X2) ---")

	// Extract odds
	var homeOdds, awayOdds string
	for _, odd := range pm.Main.SP.GameLines.Odds {
		if odd.Name == "Winner" && odd.Header == "1" {
			homeOdds = odd.Odds
		}
		if odd.Name == "Winner" && odd.Header == "2" {
			awayOdds = odd.Odds
		}
	}

	// Determine outcome
	parts := strings.Split(result.SS, "-")
	homeScore, _ := strconv.Atoi(parts[0])
	awayScore, _ := strconv.Atoi(parts[1])

	fmt.Printf("%s to Win: %s => ", result.Home.Name, homeOdds)
	if homeScore > awayScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}

	fmt.Printf("%s to Win: %s => ", result.Away.Name, awayOdds)
	if awayScore > homeScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}
	fmt.Println()
}

func evaluateCorrectScoreMarket(pm PrematchResult, result MatchResult) {
	fmt.Println("--- Correct Set Score ---")

	// Get all correct score odds
	for _, odd := range pm.Main.SP.CorrectSetScore.Odds {
		scoreType := "Home " + odd.Name
		if odd.Header == "2" {
			scoreType = "Away " + odd.Name
		}

		fmt.Printf("%s: %s => ", scoreType, odd.Odds)
		if (odd.Header == "1" && result.SS == odd.Name) ||
			(odd.Header == "2" && result.SS == odd.Name) {
			fmt.Println("WON")
		} else {
			fmt.Println("LOST")
		}
	}
	fmt.Println()
}

func evaluateTotalPointsMarket(pm PrematchResult, result MatchResult) {
	fmt.Println("--- Total Points ---")

	// Calculate total points
	totalPoints := 0
	for _, set := range result.Scores {
		home, _ := strconv.Atoi(set.Home)
		away, _ := strconv.Atoi(set.Away)
		totalPoints += home + away
	}

	// Find total points market
	for _, odd := range pm.Main.SP.GameLines.Odds {
		if odd.Name == "Total" {
			fmt.Printf("%s %s: %s (Actual: %d) => ",
				odd.Name, odd.Handicap, odd.Odds, totalPoints)

			if strings.HasPrefix(odd.Handicap, "O ") {
				var line float64
				fmt.Sscanf(odd.Handicap, "O %f", &line)
				if float64(totalPoints) > line {
					fmt.Println("WON")
				} else {
					fmt.Println("LOST")
				}
			} else if strings.HasPrefix(odd.Handicap, "U ") {
				var line float64
				fmt.Sscanf(odd.Handicap, "U %f", &line)
				if float64(totalPoints) < line {
					fmt.Println("WON")
				} else {
					fmt.Println("LOST")
				}
			}
		}
	}
	fmt.Println()
}

func evaluateHandicapMarket(pm PrematchResult, result MatchResult) {
	fmt.Println("--- Handicap ---")

	parts := strings.Split(result.SS, "-")
	homeScore, _ := strconv.Atoi(parts[0])
	awayScore, _ := strconv.Atoi(parts[1])
	scoreDiff := homeScore - awayScore

	for _, odd := range pm.Main.SP.GameLines.Odds {
		if odd.Name == "Handicap" {
			var handicap float64
			if strings.HasPrefix(odd.Handicap, "-") {
				fmt.Sscanf(odd.Handicap, "-%f", &handicap)
				fmt.Printf("%s %s: %s => ", result.Home.Name, odd.Handicap, odd.Odds)
				if float64(scoreDiff) > handicap {
					fmt.Println("WON")
				} else {
					fmt.Println("LOST")
				}
			} else if strings.HasPrefix(odd.Handicap, "+") {
				fmt.Sscanf(odd.Handicap, "+%f", &handicap)
				fmt.Printf("%s %s: %s => ", result.Away.Name, odd.Handicap, odd.Odds)
				if float64(-scoreDiff) < handicap {
					fmt.Println("WON")
				} else {
					fmt.Println("LOST")
				}
			}
		}
	}
	fmt.Println()
}

func evaluateDoubleChanceMarket(pm PrematchResult, result MatchResult) {
	fmt.Println("--- Double Chance ---")

	// Get 1X2 odds
	var homeOdds, awayOdds string
	for _, odd := range pm.Main.SP.GameLines.Odds {
		if odd.Name == "Winner" && odd.Header == "1" {
			homeOdds = odd.Odds
		}
		if odd.Name == "Winner" && odd.Header == "2" {
			awayOdds = odd.Odds
		}
	}

	parts := strings.Split(result.SS, "-")
	homeScore, _ := strconv.Atoi(parts[0])
	awayScore, _ := strconv.Atoi(parts[1])

	// 1X - Home or Draw (but volleyball rarely has draws)
	fmt.Printf("Home or Draw (1X): Derived from %s => ", homeOdds)
	if homeScore >= awayScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}

	// X2 - Draw or Away
	fmt.Printf("Draw or Away (X2): Derived from %s => ", awayOdds)
	if awayScore >= homeScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}

	// 12 - Home or Away
	fmt.Printf("Home or Away (12): Derived from %s/%s => ", homeOdds, awayOdds)
	if homeScore != awayScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}
	fmt.Println()
}

package volleyball

import (
	"bet_evaluator/models"
	"bet_evaluator/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func Evaluate() {
	// Load and parse prematch data

	prematch, err := utils.ParseData[models.PrematchData]("volleyball/volleyball_prematch.json")
	if err != nil {
		log.Fatalf("Error loading prematch data: %v", err)
	}
	// Load and parse result data
	result, err := utils.ParseData[models.ResultData]("volleyball/volleyball_result.json")
	if err != nil {
		log.Fatalf("Error loading result data: %v", err)
	}

	if len(result.Results) == 0 {
		log.Fatal("No result data found")
	}

	index := utils.GetRandomIndex(len(prematch.Results))
	index2 := utils.GetRandomIndex(len(result.Results))
	match := result.Results[index2]
	fmt.Printf("\nMatch: %s vs %s\n", match.Home.Name, match.Away.Name)
	fmt.Printf("Final Score: %s\n\n", match.SS)

	evaluateWinnerMarket(prematch.Results[index], match)
	evaluateCorrectScoreMarket(prematch.Results[index], match)
	evaluateTotalPointsMarket(prematch.Results[index], match)
	evaluateHandicapMarket(prematch.Results[index], match)
	evaluateDoubleChanceMarket(prematch.Results[index], match)

}

func evaluateWinnerMarket(pm models.PrematchResult, result models.MatchResult) {
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

func evaluateCorrectScoreMarket(pm models.PrematchResult, result models.MatchResult) {
	fmt.Println("--- Correct Set Score ---")

	// Get all correct score odds
	for _, odd := range pm.Main.SP.CorrectSetScore.Odds {
		scoreType := "Home " + odd.Name
		if odd.Header == "2" {
			scoreType = "Away " + odd.Name
		}

		fmt.Printf("%s: %s => ", scoreType, odd.Odds)
		if (odd.Header == "1" && result.SS == odd.Name) ||
			(odd.Header == "2" && result.SS == utils.ReverseScore(odd.Name)) {
			fmt.Println("WON")
		} else {
			fmt.Println("LOST")
		}
	}
	fmt.Println()
}

func evaluateTotalPointsMarket(pm models.PrematchResult, result models.MatchResult) {
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

func evaluateHandicapMarket(pm models.PrematchResult, result models.MatchResult) {
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

func evaluateDoubleChanceMarket(pm models.PrematchResult, result models.MatchResult) {
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

	// 12 - Home or Away
	fmt.Printf("Home or Away (12): Derived from %s/%s => ", homeOdds, awayOdds)
	if homeScore != awayScore {
		fmt.Println("WON")
	} else {
		fmt.Println("LOST")
	}
	fmt.Println()
}

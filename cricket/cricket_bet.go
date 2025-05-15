package cricket

import (
	"bet_evaluator/models"
	"bet_evaluator/utils"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var index int

func Evaluate() {
	// Load and parse prematch data
	prematch, err := utils.ParseData[models.CricketPrematch]("cricket/cricket_prematch.json")
	if err != nil {
		log.Fatalf("Error loading prematch data: %v", err)
	}

	// Load and parse result data
	result, err := utils.ParseData[models.CricketResult]("cricket/cricket_result.json")
	if err != nil {
		log.Fatalf("Error loading result data: %v", err)
	}

	index = utils.GetRandomIndex(len(result.Results))

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

func findMatchingResult(prematch *models.CricketPrematch, result *models.CricketResult) (*struct {
	ID          string        `json:"id"`
	SportID     string        `json:"sport_id"`
	Time        string        `json:"time"`
	TimeStatus  string        `json:"time_status"`
	League      models.League `json:"league"`
	Home        models.Team   `json:"home"`
	Away        models.Team   `json:"away"`
	SS          string        `json:"ss"`
	Extra       models.Extra  `json:"extra"`
	HasLineup   int           `json:"has_lineup"`
	ConfirmedAt string        `json:"confirmed_at"`
}, error) {
	prematchEventID := prematch.Results[index].EventID

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

func evaluateBets(prematch *models.CricketPrematch, result *struct {
	ID          string        `json:"id"`
	SportID     string        `json:"sport_id"`
	Time        string        `json:"time"`
	TimeStatus  string        `json:"time_status"`
	League      models.League `json:"league"`
	Home        models.Team   `json:"home"`
	Away        models.Team   `json:"away"`
	SS          string        `json:"ss"`
	Extra       models.Extra  `json:"extra"`
	HasLineup   int           `json:"has_lineup"`
	ConfirmedAt string        `json:"confirmed_at"`
}, homeScore, awayScore int) {
	fmt.Printf("\n=== Match Evaluation ===\n")
	fmt.Printf("League: %s\n", result.League.Name)
	fmt.Printf("Teams: %s vs %s\n", result.Home.Name, result.Away.Name)
	fmt.Printf("Venue: %s, %s\n", result.Extra.StadiumData.Name, result.Extra.StadiumData.City)
	fmt.Printf("Result: %s (%d - %d)\n", result.SS, homeScore, awayScore)

	fmt.Printf("\n--- Bet Evaluation Results ---\n")

	// Evaluate match winner market
	if len(prematch.Results[index].Main.SP.ToWinTheMatch.Odds) > 0 {
		fmt.Println("\nMatch Winner Market:")
		for _, odd := range prematch.Results[index].Main.SP.ToWinTheMatch.Odds {
			fmt.Printf("- %s @ %s => ", odd.Name, odd.Odds)
			if (odd.Name == "1" && homeScore > awayScore) || (odd.Name == "2" && awayScore > homeScore) {
				fmt.Println("WON")
			} else {
				fmt.Println("LOST")
			}
		}
	}

	// Evaluate most sixes market
	if len(prematch.Results[index].Match.SP.MostMatchSixes.Odds) > 0 {
		fmt.Println("\nMost Match Sixes Market:")
		for _, odd := range prematch.Results[index].Match.SP.MostMatchSixes.Odds {
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
	if len(prematch.Results[index].Match.SP.MostMatchFours.Odds) > 0 {
		fmt.Println("\nMost Match Fours Market:")
		for _, odd := range prematch.Results[index].Match.SP.MostMatchFours.Odds {
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

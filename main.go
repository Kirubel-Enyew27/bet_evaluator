package main

import (
	"bet_evaluator/cricket"
	"bet_evaluator/volleyball"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		// Prompt user for input
		fmt.Print("Enter your choice from sport betting [volleyball, cricket] or 'exit' to quit: ")

		var sportBet string
		if _, err := fmt.Scanln(&sportBet); err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Normalize input (trim spaces and lowercase)
		sportBet = strings.TrimSpace(strings.ToLower(sportBet))

		// Process input
		switch sportBet {
		case "volleyball":
			volleyball.Evaluate()
			return
		case "cricket":
			cricket.Evaluate()
			return
		case "exit":
			fmt.Println("Exiting program...")
			os.Exit(0)
		default:
			fmt.Printf("Invalid choice '%s'. Please try again.\n", sportBet)
		}
	}
}

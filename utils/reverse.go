package utils

import "strings"

func ReverseScore(score string) string {
	parts := strings.Split(score, "-")
	if len(parts) != 2 {
		return score
	}
	return parts[1] + "-" + parts[0]
}

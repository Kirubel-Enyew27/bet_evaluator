package utils

import (
	"bet_evaluator/models"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

func ParseData[T any](filename string) (*T, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var matchData T
	if err := json.Unmarshal(data, &matchData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	// Get the concrete type name for better error messages
	typeName := fmt.Sprintf("%T", matchData)

	// Type-specific validation
	switch v := any(&matchData).(type) {
	case *models.PrematchData:
		if len(v.Results) == 0 {
			return nil, fmt.Errorf("invalid volleyball prematch data")
		}
	case *models.ResultData:
		if len(v.Results) == 0 {
			return nil, fmt.Errorf("invalid volleyball result data")
		}
	case *models.CricketPrematch:
		if v.Success != 1 || len(v.Results) == 0 {
			return nil, fmt.Errorf("invalid cricket prematch data: success=%d, results=%d", v.Success, len(v.Results))
		}

	case *models.CricketResult:
		if v.Success != 1 || len(v.Results) == 0 {
			return nil, fmt.Errorf("invalid cricket result data: success=%d, results=%d", v.Success, len(v.Results))
		}
	default:
		if reflect.ValueOf(matchData).IsNil() {
			return nil, fmt.Errorf("invalid prematch data: nil value for type %s", typeName)
		}
	}

	return &matchData, nil
}

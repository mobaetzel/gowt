package utils

import "encoding/json"

type Values map[string]any

func ParseValues(values string) (Values, error) {
	var data Values
	err := json.Unmarshal([]byte(values), &data)
	return data, err
}

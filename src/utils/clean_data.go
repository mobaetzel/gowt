package utils

import "html"

func CleanData(data any) any {
	if mapData, ok := data.(map[string]any); ok {
		for k, v := range mapData {
			mapData[k] = CleanData(v)
		}
	} else if sliceData, ok := data.([]any); ok {
		for i, v := range sliceData {
			sliceData[i] = CleanData(v)
		}
	} else if strData, ok := data.(string); ok {
		return html.EscapeString(strData)
	}
	return data
}

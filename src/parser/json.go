package parser

import (
	"encoding/json"
)

func JSON(body []byte) map[string]any {

	var jsonbody map[string]any

	err := json.Unmarshal(body, &jsonbody)

	if err != nil {
		return nil
	}

	return jsonbody
}

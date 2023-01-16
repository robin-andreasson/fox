package parser

import (
	"encoding/json"
)

func JSON(body []byte) any {

	var jsonbody interface{}

	err := json.Unmarshal(body, &jsonbody)

	if err != nil {
		return nil
	}

	return jsonbody
}

package services

import (
	"encoding/json"
)

func BuildRequestBody(limit int, searchAfter []interface{}) (string, error) {
	body := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"@timestamp": map[string]string{
					"order": "desc",
				},
				"_id": map[string]string{
					"order": "desc",
				},
			},
		},
		"size": limit,
	}

	if len(searchAfter) > 0 {
		body["search_after"] = searchAfter
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(jsonBody), nil
}

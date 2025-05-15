package dto

type WazuhResponse struct {
	Hits struct {
		Hits []struct {
			ID     string `json:"_id"`
			Source struct {
				Agent struct {
					IP   string `json:"ip"`
					Name string `json:"name"`
				} `json:"agent"`
				Data      map[string]interface{} `json:"data"`
				Timestamp string                 `json:"time"`
				Time      string                 `json:"@timestamp"`
				Location  string                 `json:"location"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type SimplifiedLog struct {
	ID        string                 `json:"id"`
	AgentIP   string                 `json:"agent_ip"`
	AgentName string                 `json:"agent_name"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"time"`
	Time      string                 `json:"@timestamp"`
	Location  string                 `json:"location"`
}

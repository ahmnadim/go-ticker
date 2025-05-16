package repositories

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"poc/services"
	dto "poc/structs"
	"strconv"
)

func FetchFromWazuh(offset int, search_after []interface{}) ([]dto.SimplifiedLog, []interface{}, error) {
	limit, err := strconv.Atoi(os.Getenv("DATA_LIMIT"))
	if err != nil {
		limit = 10
	}

	var res dto.WazuhResponse

	wazuh_url := os.Getenv("WAZUH_API_URL")
	wazuh_username := os.Getenv("WAZUH_USERNAME")
	wazuh_password := os.Getenv("WAZUH_PASSWORD")

	// Define the JSON body (same as the cURL -d part)
	jsonBody, _ := services.BuildRequestBody(limit, search_after)
	// Convert the body to JSON
	// fmt.Println("body: ", jsonBody, wazuh_url)

	req, err := http.NewRequest("GET", wazuh_url, bytes.NewBuffer([]byte(jsonBody)))
	if err != nil {
		panic(err)
	}

	// Encode credentials to base64 for Basic Auth
	auth := wazuh_username + ":" + wazuh_password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	// client := &http.Client{}
	//TODO: need to fix
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ⚠️ NOT safe for production
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read and print the response body
	res_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(res_body, &res); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	var logs []dto.SimplifiedLog
	var sort []interface{}
	for _, hit := range res.Hits.Hits {
		logEntry := dto.SimplifiedLog{
			ID:        hit.ID,
			AgentIP:   hit.Source.Agent.IP,
			AgentName: hit.Source.Agent.Name,
			Data:      hit.Source.Data,
			Sort:      hit.Sort,
		}
		logs = append(logs, logEntry)
		sort = logEntry.Sort
	}
	if len(logs) == 0 {
		return nil, nil, err
	}

	return logs, sort, nil

}

package repositories

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	dto "poc/structs"
	"strconv"
)

func FetchFromWazuh(offset int) ([]dto.SimplifiedLog, error) {
	limit, err := strconv.Atoi(os.Getenv("DATA_LIMIT"))
	if err != nil {
		limit = 10
	}

	var res dto.WazuhResponse

	wazuh_url := os.Getenv("WAZUH_API_URL")
	wazuh_username := os.Getenv("WAZUH_USERNAME")
	wazuh_password := os.Getenv("WAZUH_PASSWORD")

	// Parse the base URL
	baseURL, err := url.Parse(wazuh_url)
	if err != nil {
		panic(err)
	}

	// Add query parameters
	params := url.Values{}
	params.Add("size", strconv.Itoa(limit))
	params.Add("from", strconv.Itoa(offset))

	// Attach the params to the base URL
	baseURL.RawQuery = params.Encode()
	fmt.Println("base url: ", baseURL.String())
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		panic(err)
	}

	// Encode credentials to base64 for Basic Auth
	auth := wazuh_username + ":" + wazuh_password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encodedAuth)

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	var logs []dto.SimplifiedLog
	for _, hit := range res.Hits.Hits {
		logEntry := dto.SimplifiedLog{
			ID:        hit.ID,
			AgentIP:   hit.Source.Agent.IP,
			AgentName: hit.Source.Agent.Name,
			Data:      hit.Source.Data,
		}
		logs = append(logs, logEntry)
	}
	if len(logs) == 0 {
		return nil, err
	}

	return logs, nil

}

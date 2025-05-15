package main

import (
	"fmt"
	"os"
	"os/signal"
	"poc/initializers"
	"poc/models"
	"poc/repositories"
	"strconv"
	"syscall"
	"time"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
}

// Struct to hold data in memory
type Data struct {
	Previous    []models.User
	APIResponse []models.User
	Timestamp   string
}

// Global in-memory storage
var memoryData Data
var offset = 0

func main() {

	fmt.Println("Press CTRL + c to stop the ticker.")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	stopChan := make(chan bool)

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		stopChan <- true
	}()

	for {
		select {
		case <-stopChan:
			fmt.Println("Stopping ticker.")
			return

		case <-ticker.C:
			func() {

				start := time.Now()

				// Fetch new API data
				apiData, err := repositories.FetchDataFromDB(offset)
				fmt.Println("Api response: ", offset, len(apiData) == 0, len(apiData), apiData)

				if len(apiData) == 0 {
					fmt.Println("Data fetching completed, Stopping ticker.")
					stopChan <- true
					return
				}
				if err != nil {
					fmt.Println("Error fetching from API:", err)
					return
				}

				// Show current state
				// fmt.Printf("Tick:\nPrevious: %d\nAPI Response: %d\n---\n", memoryData.APIResponse, apiData)
				// fmt.Println("Tick Previous: ", memoryData.APIResponse)

				// Update memoryData struct
				memoryData = Data{
					Previous:    memoryData.APIResponse,
					APIResponse: apiData,
					Timestamp:   time.Now().Format(time.RFC3339),
				}

				fmt.Printf("Tick done in: %v\n\n", time.Since(start))
				limit, err := strconv.Atoi(os.Getenv("DATA_LIMIT"))
				if err != nil {
					limit = 10
				}
				offset += limit
			}()
		}
	}

}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"poc/initializers"
	"poc/repositories"
	"strconv"
	"syscall"
	"time"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectDB()
}

var offset = 0
var sort []interface{}

func main() {
	fmt.Println("Press CTRL + c to stop the ticker.")
	tickerInterval, err := strconv.Atoi(os.Getenv("TICKER_INTERVAL"))
	if err != nil {
		tickerInterval = 30
	}
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
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

				logs, _sort, err := repositories.FetchFromWazuh(offset, sort)
				fmt.Println("longs: ", len(logs), sort)
				sort = _sort
				if len(logs) == 0 {
					fmt.Println("Data fetching completed, Stopping ticker.")
					stopChan <- true
					return
				}
				if err != nil {
					fmt.Println("Error fetching from API:", err)
					return
				}

				fmt.Printf("Tick done in: %v\n\n", time.Since(start))
				limit, err := strconv.Atoi(os.Getenv("DATA_LIMIT"))
				if err != nil {
					limit = 10
				}
				offset += 9900 + limit
			}()
		}
	}

}

package modulyn

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var datastore store

func Initialize(environmentID string, applicationID string) error {
	if environmentID == "" {
		return fmt.Errorf("environmentID cannot be empty")
	}

	if applicationID == "" {
		applicationID = uuid.NewString()
	}

	doneChan := make(chan bool, 1)

	go func() {
		sseURL := fmt.Sprintf("http://localhost:8080/events?sdk_key=%s&appid=%s", environmentID, applicationID)
		response, err := http.Get(sseURL)
		if err != nil {
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return
		}

		fmt.Printf("Successfully connected to modulyn stream\n")

		reader := bufio.NewReader(response.Body)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			line = strings.TrimSpace(line)

			if after, ok := strings.CutPrefix(line, "data:"); ok {
				line = after
				line = strings.TrimSpace(line)
				fmt.Printf("Received data: %s\n", line)

				var event Event
				if err := json.Unmarshal([]byte(line), &event); err != nil {
					return
				}

				fmt.Printf("Event Type: %s", event.Type)

				var features []Feature
				if event.Type == "all_features" {
					if err := json.Unmarshal(event.Data, &features); err != nil {
						return
					}

					for _, feature := range features {
						datastore.addOrUpdate(feature)
						fmt.Printf("Feature ID: %s, Name: %s, Enabled: %t\n", feature.ID, feature.Name, feature.Enabled)
					}

					doneChan <- true
				}
			}
		}
	}()

	for range 1 {
		isDone := <-doneChan
		if !isDone {
			return fmt.Errorf("error initializing modulyn server")
		}
	}

	return nil
}

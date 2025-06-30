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

				switch event.Type {
				case "all_features":
					var features []Feature
					if err := json.Unmarshal(event.Data, &features); err != nil {
						return
					}

					for _, feature := range features {
						datastore.addOrUpdate(feature)
						fmt.Printf("Feature ID: %s, Name: %s, Enabled: %t\n", feature.ID, feature.Name, feature.Enabled)
					}

					doneChan <- true
				case "feature_created", "feature_updated":
					var newFeature Feature
					if err := json.Unmarshal(event.Data, &newFeature); err != nil {
						return
					}

					datastore.addOrUpdate(newFeature)
					fmt.Printf("Feature Created - ID: %s, Name: %s, Enabled: %t\n", newFeature.ID, newFeature.Name, newFeature.Enabled)

				case "feature_deleted":
					var deletedFeature Feature
					if err := json.Unmarshal(event.Data, &deletedFeature); err != nil {
						return
					}

					datastore.remove(deletedFeature)
					fmt.Printf("Feature Deleted - ID: %s, Name: %s\n", deletedFeature.ID, deletedFeature.Name)
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

package modulyn

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

	sseURL := fmt.Sprintf("http://localhost:8080/events?sdk_key=%s&appid=%s", environmentID, applicationID)
	retryDelay := 5 * time.Second

	go func() {
		for {
			response, err := http.Get(sseURL)
			if err != nil {
				log.Println("error connecting to modulyn stream:", err)
				time.Sleep(retryDelay)
				continue
			}

			if response.StatusCode != http.StatusOK {
				log.Println("received non-200 response code from modulyn stream:", response.StatusCode)
				response.Body.Close()
				time.Sleep(retryDelay)
				continue
			}

			fmt.Printf("Successfully connected to modulyn stream\n")

			reader := bufio.NewReader(response.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					log.Println("error reading from modulyn stream:", err)
					response.Body.Close()
					break // triggers reconnect
				}

				line = strings.TrimSpace(line)

				if after, ok := strings.CutPrefix(line, "data:"); ok {
					line = after
					line = strings.TrimSpace(line)

					var event Event
					if err := json.Unmarshal([]byte(line), &event); err != nil {
						log.Println("error unmarshalling event from modulyn stream:", err)
						response.Body.Close()
						break // triggers reconnect
					}

					switch event.Type {
					case "all_features":
						var features []Feature
						if err := json.Unmarshal(event.Data, &features); err != nil {
							log.Println("error unmarshalling features from modulyn stream:", err)
							response.Body.Close()
							break
						}
						for _, feature := range features {
							datastore.addOrUpdate(feature)
						}
					case "feature_created", "feature_updated":
						var newFeature Feature
						if err := json.Unmarshal(event.Data, &newFeature); err != nil {
							log.Println("error unmarshalling feature from modulyn stream:", err)
							response.Body.Close()
							break
						}
						datastore.addOrUpdate(newFeature)
					case "feature_deleted":
						var deletedFeature Feature
						if err := json.Unmarshal(event.Data, &deletedFeature); err != nil {
							log.Println("error unmarshalling feature from modulyn stream:", err)
							response.Body.Close()
							break
						}
						datastore.remove(deletedFeature)
					}
				}
			}
			log.Println("Disconnected from modulyn stream, retrying in 5 seconds...")
			time.Sleep(retryDelay)
		}
	}()

	return nil
}

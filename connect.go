package modulyn

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func Initialize(environmentID string, applicationID string) error {
	if environmentID == "" {
		return fmt.Errorf("environmentID cannot be empty")
	}

	if applicationID == "" {
		applicationID = uuid.NewString()
	}

	sseURL := fmt.Sprintf("http://localhost:8080/events?sdk_key=%s&appid=%s", environmentID, applicationID)
	response, err := http.Get(sseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to SSE server: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to connect to SSE server, status code: %d", response.StatusCode)
	}

	fmt.Printf("Successfully connected to modulyn stream\n")

	reader := bufio.NewReader(response.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading from SSE stream: %w", err)
		}

		line = strings.TrimSpace(line)

		if after, ok := strings.CutPrefix(line, "data:"); ok {
			line = after
			line = strings.TrimSpace(line)
			fmt.Printf("Received data: %s\n", line)

			var event Event
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				return fmt.Errorf("error unmarshalling event data: %w", err)
			}

			fmt.Printf("Event Type: %s", event.Type)

			var features []Feature
			if event.Type == "all_features" {
				if err := json.Unmarshal(event.Data, &features); err != nil {
					return fmt.Errorf("error unmarshalling features: %w", err)
				}

				for _, feature := range features {
					fmt.Printf("Feature ID: %s, Name: %s, Enabled: %t\n", feature.ID, feature.Name, feature.Enabled)
				}
			}
		}
	}
}

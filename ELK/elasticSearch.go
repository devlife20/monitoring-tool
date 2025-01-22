package ELK

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"io"
	"time"
)

// FetchLogs retrieves and prints all documents from the specified Elasticsearch index
func FetchLogs(index string) error {
	// Create client
	cfg := elasticsearch.Config{
		Addresses: []string{
			viper.GetString("elastic_url"),
		},
		APIKey: viper.GetString("elastic_api_key"),
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	// Variables for pagination and tailing
	var (
		from      = 0
		pageSize  = 100 // Number of documents per request
		totalHits int
		lastTime  time.Time // Track the timestamp of the last document fetched
	)

	// Initialize lastTime to the current time
	lastTime = time.Now()

	for {
		// Perform the search request with pagination
		res, err := es.Search(
			es.Search.WithContext(context.Background()),
			es.Search.WithIndex(index),
			es.Search.WithFrom(from),
			es.Search.WithSize(pageSize),
			es.Search.WithTrackTotalHits(true),    // Track total hits
			es.Search.WithSort("@timestamp:desc"), // Sort by timestamp to get the latest documents first
			es.Search.WithPretty(),
		)
		if err != nil {
			return fmt.Errorf("error searching: %w", err)
		}
		defer res.Body.Close()

		// Check for response errors
		if res.IsError() {
			body, _ := io.ReadAll(res.Body)
			return fmt.Errorf("error response from Elasticsearch: %s", body)
		}

		// Parse the response body
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return fmt.Errorf("error parsing response: %w", err)
		}

		// Extract and display hits
		hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
		if totalHits == 0 {
			totalHits = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
			fmt.Printf("Total documents: %d\n", totalHits)
		}

		newDocs := 0
		for _, hit := range hits {
			// Retrieve the full document (_source)
			doc := hit.(map[string]interface{})["_source"]

			// Check if the document is newer than the last fetched document
			timestampStr := doc.(map[string]interface{})["@timestamp"].(string)
			timestamp, err := time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				fmt.Printf("Error parsing timestamp: %v\n", err)
				continue
			}

			if timestamp.After(lastTime) {
				// Convert document to JSON string for pretty-printing
				docJSON, _ := json.MarshalIndent(doc, "", "  ")
				fmt.Println(string(docJSON))

				// Update lastTime to the timestamp of the latest document
				if timestamp.After(lastTime) {
					lastTime = timestamp
				}
				newDocs++
			}
		}

		// If no new documents were found, wait for a while before checking again
		if newDocs == 0 {
			fmt.Println("No new documents. Waiting for new data...")
			time.Sleep(10 * time.Second) // Adjust the sleep duration as needed
		}

		// Move to the next page
		from += pageSize
	}
}

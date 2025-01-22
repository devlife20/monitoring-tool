package ELK

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
	"time"
)

func WatchElasticLogs(index string) error {
	printed := false
	const interval = 3 * time.Second
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Fatalf("failed to read config: %v", err)
		}
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			viper.GetString("elk.elastic_url"),
		},
		APIKey: viper.GetString("elk.api_token"),
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Use a map to track seen document IDs
	seenDocs := make(map[string]bool)
	var mu sync.Mutex // Mutex to safely access seenDocs

	fmt.Printf(" Fetching logs from index '%s'\n \n", index)

	for {
		select {
		case <-ticker.C:
			query := `{
                "query": {
                    "match_all": {}
                }
            }`

			res, err := es.Search(
				es.Search.WithContext(context.Background()),
				es.Search.WithIndex(index),
				es.Search.WithBody(strings.NewReader(query)),
				es.Search.WithSize(100),
				es.Search.WithTrackTotalHits(true),
			)
			if err != nil {
				fmt.Printf("Search error: %v\n", err)
				continue
			}

			var result map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				res.Body.Close()
				continue
			}
			res.Body.Close()

			if hits, ok := result["hits"].(map[string]interface{}); ok {
				if total, ok := hits["total"].(map[string]interface{}); ok {
					totalDocs := int(total["value"].(float64))
					if totalDocs == 0 {
						continue
					}
				}

				if hitsList, ok := hits["hits"].([]interface{}); ok {
					newDocs := false
					for _, hit := range hitsList {
						if doc, ok := hit.(map[string]interface{}); ok {
							// Get document ID
							docID := doc["_id"].(string)

							// Check if we've seen this document before
							mu.Lock()
							seen := seenDocs[docID]
							if !seen {
								seenDocs[docID] = true
								newDocs = true
								mu.Unlock()

								if source, ok := doc["_source"]; ok {
									output, err := json.MarshalIndent(source, "", "  ")
									if err != nil {
										continue
									}
									fmt.Printf("log (ID: %s):\n%s\n", docID, output)
								}
							} else {
								mu.Unlock()
							}
						}
					}

					if !newDocs && !printed {
						fmt.Printf("Watching for new logs....\n")
						printed = true
					} else if newDocs {
						printed = false
					}

				}
			}
		}
	}
}

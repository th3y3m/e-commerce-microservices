package elasticsearch_server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/spf13/viper"
)

func ConnectToElasticsearch() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			viper.GetString("ELASTICSEARCH_URL"), // Update with your Elasticsearch server URL
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	// Check if the connection is successful
	for i := 0; i < 5; i++ { // Retry up to 5 times
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		res, err := es.Ping(es.Ping.WithContext(ctx))
		if err == nil && !res.IsError() {
			log.Println("Successfully connected to Elasticsearch")
			return es, nil
		}

		log.Printf("Failed to connect to Elasticsearch (attempt %d/5): %v", i+1, err)
		time.Sleep(1 * time.Second) // Wait before retrying
	}

	return nil, fmt.Errorf("failed to connect to Elasticsearch after multiple attempts: %w", err)
}

func CreateIndex(es *elasticsearch.Client, indexName string) error {
	// Check if the index already exists
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("Index %s already exists", indexName)
		return nil
	}

	// if res.StatusCode == 200 {
	// 	// Delete the existing index
	// 	res, err := es.Indices.Delete([]string{indexName})
	// 	if err != nil {
	// 		return fmt.Errorf("error deleting existing index: %w", err)
	// 	}
	// 	defer res.Body.Close()
	// 	if res.IsError() {
	// 		bodyBytes, _ := io.ReadAll(res.Body)
	// 		return fmt.Errorf("error response from Elasticsearch while deleting index: %s - %s", res.Status(), string(bodyBytes))
	// 	}
	// 	log.Printf("Index %s deleted to apply new mappings", indexName)
	// }

	// Define index settings and mappings
	indexSettings := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"filter": map[string]interface{}{
					"my_phonetic_filter": map[string]interface{}{
						"type":    "phonetic",
						"encoder": "doublemetaphone",
						"replace": false,
					},
				},
				"analyzer": map[string]interface{}{
					"my_phonetic_analyzer": map[string]interface{}{
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"my_phonetic_filter",
						},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"ProductName": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"phonetic": map[string]interface{}{ // Phonetic subfield
							"type":     "text",
							"analyzer": "my_phonetic_analyzer",
						},
						"keyword": map[string]interface{}{ // Exact match subfield
							"type": "keyword",
						},
					},
				},
			},
		},
	}

	// Convert settings to JSON
	settingsJSON, err := json.Marshal(indexSettings)
	if err != nil {
		return fmt.Errorf("error marshaling index settings: %w", err)
	}

	// Create the index
	res, err = es.Indices.Create(indexName, es.Indices.Create.WithBody(bytes.NewReader(settingsJSON)))
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("error response from Elasticsearch: %s - %s", res.Status(), string(bodyBytes))
	}

	log.Printf("Index %s created successfully", indexName)
	return nil
}

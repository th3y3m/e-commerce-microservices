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
			"index": map[string]interface{}{
				"max_ngram_diff": 17, // Allow for a larger n-gram difference
			},
			"analysis": map[string]interface{}{
				"filter": map[string]interface{}{
					"my_phonetic_filter": map[string]interface{}{
						"type":    "phonetic",
						"encoder": "doublemetaphone",
						"replace": false,
					},
					"synonym_filter_1": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							"laptop, notebook, portable computer",
						},
					},
					"synonym_filter_2": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							"nokia, mobile phone, cellphone",
						},
					},
					"synonym_filter_3": map[string]interface{}{
						"type": "synonym",
						"synonyms": []string{
							"computer, pc, desktop",
						},
					},
					// "synonym_filter": map[string]interface{}{
					// 	"type": "synonym",
					// 	"synonyms": []string{
					// 		"café, cafe, coffee shop",
					// 		"laptop, notebook, portable computer",
					// 		"nokia, mobile phone, cellphone",
					// 		"computer, pc, desktop",
					// 		"clothes, clothing, apparel",
					// 		"book, novel, paperback",
					// 		"bag, backpack, purse, handbag",
					// 		"shoes, sneakers, footwear, boots",
					// 		"iphone, smartphone, mobile phone",
					// 		"shirt, top, blouse",
					// 		"t-shirt, tee, tee shirt",
					// 	},
					// },
					"word_delimiter_custom": map[string]interface{}{
						"type":                  "word_delimiter_graph",
						"generate_word_parts":   true,
						"generate_number_parts": true,
						"catenate_words":        true,
						"catenate_numbers":      true,
						"catenate_all":          true,
						"preserve_original":     true,
						"split_on_case_change":  false,
					},
					"ngram_filter": map[string]interface{}{
						"type":     "ngram",
						"min_gram": 3,
						"max_gram": 20,
					},
					"my_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "english",
					},
					"shingle_filter": map[string]interface{}{
						"type":             "shingle",
						"min_shingle_size": 2,
						"max_shingle_size": 3,
					},
				},
				"analyzer": map[string]interface{}{
					"my_custom_analyzer": map[string]interface{}{
						"tokenizer": "whitespace",
						"filter": []string{
							"lowercase",
							"synonym_filter_1",
							"synonym_filter_2",
							"synonym_filter_3",
							"word_delimiter_custom",
							"my_phonetic_filter",
							"ngram_filter",
							"my_stemmer",
							"shingle_filter",
						},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"ProductName": map[string]interface{}{
					"type":     "text",
					"analyzer": "my_custom_analyzer",
					"fields": map[string]interface{}{
						"raw": map[string]interface{}{
							"type": "keyword",
						},
						"folded": map[string]interface{}{
							"type":     "text",
							"analyzer": "my_custom_analyzer",
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

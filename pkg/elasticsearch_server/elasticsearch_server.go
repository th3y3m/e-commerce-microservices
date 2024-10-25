package elasticsearch_server

import (
	"context"
	"fmt"
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

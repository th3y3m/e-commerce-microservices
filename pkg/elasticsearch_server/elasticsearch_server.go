package elasticsearch_server

import (
	"log"

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
		return nil, err
	}

	log.Println("Connected to Elasticsearch")

	return es, nil
}

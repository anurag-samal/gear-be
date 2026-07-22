package search

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github/anurag/altar-be/system/config"
)

func NewElasticSearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Elasticsearch.URL},
	})
	if err != nil {
		return nil, fmt.Errorf("elasticsearch client: %w", err)
	}

	res, err := esClient.Ping()
	if err != nil {
		return nil, fmt.Errorf("elasticsearch ping: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch ping: %s", res.String())
	}

	return esClient, nil
}

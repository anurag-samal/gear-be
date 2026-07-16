package search

import (
	"github/anurag/altar-be/system/config"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
	elasticConfig := elasticsearch.Config{
		Addresses: []string{
			cfg.ES_URL,
		},
	}

	esClient, err := elasticsearch.NewClient(elasticConfig)
	if err != nil {
		return nil, err
	}

	return esClient, nil
}

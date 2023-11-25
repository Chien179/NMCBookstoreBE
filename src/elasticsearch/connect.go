package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

func GetElasticSearch() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	return es, err
}

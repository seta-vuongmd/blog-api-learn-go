package config

import (
	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func ConnectElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://elasticsearch:9200"},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	ESClient = client
}

package config

import (
	"fmt"
	"strings"

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

	// Check if posts index exists, if not create it
	res, err := client.Indices.Exists([]string{"posts"})
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// Create posts index with mapping for tags
		mapping := `{
			"mappings": {
				"properties": {
					"id": { "type": "long" },
					"title": { "type": "text" },
					"content": { "type": "text" },
					"tags": { "type": "text" }
				}
			}
		}`

		res, err := client.Indices.Create(
			"posts",
			client.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		if res.IsError() {
			fmt.Println("Error creating index:", res.String())
		}
	}
}

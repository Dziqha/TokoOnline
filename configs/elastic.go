package configs

import (
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/responses"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

var ESClient *elasticsearch.Client
const SearchIndex = "product"

func ESClientConnection() {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ESClient = client
}

func ESCreateIndexIfNotExist() {
	// Periksa apakah indeks ada
	exists, err := esapi.IndicesExistsRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)
	if err != nil {
		log.Fatalf("Error checking index existence: %s", err)
	}

	if exists.StatusCode == 404 {
		// Jika indeks tidak ada, buat indeks
		mapping := `{
			"mappings": {
				"properties": {
					"suggest": {
						"type": "completion"
					}
				}
			}
		}`
		req := esapi.IndicesCreateRequest{
			Index: SearchIndex,
			Body:  strings.NewReader(mapping),
		}

		res, err := req.Do(context.Background(), ESClient)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Fatalf("Error response from Elasticsearch: %s", res.String())
		} else {
			fmt.Println("Index created successfully.")
		}
	} else if exists.StatusCode == 200 {
		fmt.Println("Index already exists.")
	} else {
		log.Fatalf("Unexpected response status code: %d", exists.StatusCode)
	}
}

func AddProductToIndex(product models.Product) error {
	suggest := map[string]interface{}{
		"suggest": map[string]interface{}{
			"input": []string{product.Name},
		},
	}
	productJson, err := json.Marshal(suggest)
	if err != nil {
		return fmt.Errorf("error marshalling product: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      SearchIndex,
		Body:       bytes.NewReader(productJson),
		Refresh:    "true",
		DocumentID: product.ID,
	}

	res, err := req.Do(context.Background(), ESClient)
	if err != nil {
		return fmt.Errorf("error indexing product: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response from Elasticsearch: %s", res.String())
	}

	return nil
}
func UpdatedProductToIndex( esClient *elasticsearch.Client,productID string, updateData responses.Product) error {
	productJson, err := json.Marshal(updateData)
	if err != nil {
		log.Fatalf("Error marshalling product: %s", err)
	}

	req := esapi.IndexRequest{
		Index: SearchIndex,
		Body:  strings.NewReader(string(productJson)),
		Refresh: "true",
		DocumentID: productID,
	}

	res, err := req.Do(context.Background(),ESClient)
	if err != nil {
		log.Fatalf("Error updating product: %s", err)
	}
	defer res.Body.Close()

	if res.IsError(){
		log.Fatalf("Error response from Elasticsearch: %s", res.String())
	}


	return nil
}

func DeleteProductFromIndex(id string) error {
	req := esapi.DeleteRequest{
		Index: SearchIndex,
		DocumentID: id,
	}

	res, err := req.Do(context.Background(),ESClient)
	if err != nil {
		log.Fatalf("Error deleting product: %s", err)
	}
	defer res.Body.Close()

	if res.IsError(){
		log.Fatalf("Error response from Elasticsearch: %s", res.String())
	}

	return nil
}

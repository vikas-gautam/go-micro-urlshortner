package elasticsearch

import (
	"context"
	"fmt"
	"log"
	"strings"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/sirupsen/logrus"
)

var elasticClient *elastic.Client

func ConnectionElastic(conn *elastic.Client) {
	elasticClient = conn
}

func IndexingData() {

	// Index a document
	doc := `{"title": "Go and Elasticsearch", "content": "A tutorial on how to use Go and Elasticsearch together"}`
	req := esapi.IndexRequest{
		Index:      "articles",
		DocumentID: "1",
		Body:       strings.NewReader(doc),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), elasticClient)
	if err != nil {
		logrus.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
}

func SearchingData() {

	// Search for documents
	query := `{"query": {"match": {"title": "Go"}}}`
	req := esapi.SearchRequest{
		Index: []string{"articles"},
		Body:  strings.NewReader(query),
	}

	res, err := req.Do(context.Background(), elasticClient)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()

	fmt.Println(res)
}

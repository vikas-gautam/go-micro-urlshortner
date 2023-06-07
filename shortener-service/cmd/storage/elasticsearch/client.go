package elasticsearch

import (
	"errors"
	"fmt"
	"os"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
)

func ConnectToElasticsearch() (*elastic.Client, error) {
	elasticEndpoint := os.Getenv("ELASTIC_ENDPOINT")

	var (
		clusterURLs = []string{elasticEndpoint}
	)

	esConfig := elastic.Config{
		Addresses: clusterURLs,
	}
	esClient, err := elastic.NewClient(esConfig)
	if err != nil {
		return esClient, errors.New(fmt.Sprintf("connectToElasticsearch: elastic.NewClient %s", err))
	}
	res, err := esClient.Info()

	if err != nil {
		logrus.Error("Error getting response: %s", err)
		return esClient, err
	}

	defer res.Body.Close()

	logrus.Info("Connected to elasticsearch")
	return esClient, nil
}

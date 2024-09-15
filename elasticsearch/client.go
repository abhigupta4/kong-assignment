package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kong/config"
	"net/url"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

type ElasticSearchClient struct {
	Logger *zap.Logger
	client *opensearch.Client
}

func NewElasticSearchClient(config config.ElasticSearchConfig, logger *zap.Logger) (*ElasticSearchClient, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{config.Address},
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchClient{
		client: client,
		Logger: logger,
	}, nil
}

func (es *ElasticSearchClient) Write(record interface{}, index string, docID string) error {
	jsonBytes, err := json.Marshal(record)
	if err != nil {
		es.Logger.Error("Failed to marshal record to JSON", zap.Any("error", err))
		return fmt.Errorf("failed to marshal record to JSON: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: url.PathEscape(docID),
		Body:       bytes.NewReader(jsonBytes),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		es.Logger.Error("Failed to execute OpenSearch request", zap.Any("error", err))
		return fmt.Errorf("failed to execute OpenSearch request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		es.Logger.Error("Error indexing document", zap.String("docID", docID), zap.String("response", res.String()))
		return fmt.Errorf("error indexing document ID=%s: %s", docID, res.String())
	}

	es.Logger.Info("Successfully wrote document to OpenSearch", zap.String("docID", docID))
	return nil
}

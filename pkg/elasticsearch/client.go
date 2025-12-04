package elasticsearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client wraps Elasticsearch client
type Client struct {
	es      *elasticsearch.Client
	indices map[string]string
}

// Config holds Elasticsearch configuration
type Config struct {
	Addresses []string
	Username  string
	Password  string
	CloudID   string
	APIKey    string
	TLS       *TLSConfig
}

type TLSConfig struct {
	InsecureSkipVerify bool
	CertFile           string
	KeyFile            string
}

// NewClient creates a new Elasticsearch client
func NewClient(cfg Config) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CloudID:   cfg.CloudID,
		APIKey:    cfg.APIKey,
	}

	// Configure TLS if provided
	if cfg.TLS != nil {
		esCfg.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.TLS.InsecureSkipVerify,
			},
		}
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	// Test connection
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting cluster info: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from elasticsearch: %s", res.String())
	}

	return &Client{
		es:      es,
		indices: make(map[string]string),
	}, nil
}

// CreateIndex creates a new index with mapping
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping map[string]interface{}) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("error marshaling mapping: %w", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  esapi.NewJSONReader(body),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response creating index: %s", res.String())
	}

	c.indices[indexName] = indexName
	return nil
}

// DeleteIndex deletes an index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	req := esapi.IndicesDeleteRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error deleting index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error response deleting index: %s", res.String())
	}

	delete(c.indices, indexName)
	return nil
}

// IndexExists checks if index exists
func (c *Client) IndexExists(ctx context.Context, indexName string) (bool, error) {
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return false, fmt.Errorf("error checking index existence: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// IndexDocument indexes a document
func (c *Client) IndexDocument(ctx context.Context, indexName, documentID string, document interface{}) error {
	body, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       esapi.NewJSONReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response indexing document: %s", res.String())
	}

	return nil
}

// UpdateDocument updates a document
func (c *Client) UpdateDocument(ctx context.Context, indexName, documentID string, document interface{}) error {
	doc := map[string]interface{}{
		"doc": document,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       esapi.NewJSONReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response updating document: %s", res.String())
	}

	return nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(ctx context.Context, indexName, documentID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: documentID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error response deleting document: %s", res.String())
	}

	return nil
}

// BulkIndexDocuments indexes multiple documents in bulk
func (c *Client) BulkIndexDocuments(ctx context.Context, indexName string, documents []BulkDocument) error {
	if len(documents) == 0 {
		return nil
	}

	var bulkBody []byte
	for _, doc := range documents {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": indexName,
				"_id":    doc.ID,
			},
		}
		metaBytes, _ := json.Marshal(meta)
		docBytes, _ := json.Marshal(doc.Document)

		bulkBody = append(bulkBody, metaBytes...)
		bulkBody = append(bulkBody, '\n')
		bulkBody = append(bulkBody, docBytes...)
		bulkBody = append(bulkBody, '\n')
	}

	req := esapi.BulkRequest{
		Body:    esapi.NewJSONReader(bulkBody),
		Refresh: "true",
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error bulk indexing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response bulk indexing: %s", res.String())
	}

	// Parse response to check for errors
	var bulkRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkRes); err != nil {
		return fmt.Errorf("error parsing bulk response: %w", err)
	}

	if errors, ok := bulkRes["errors"].(bool); ok && errors {
		return fmt.Errorf("bulk indexing had errors")
	}

	return nil
}

// Search performs a search query
func (c *Client) Search(ctx context.Context, indexName string, query map[string]interface{}) (*SearchResponse, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  esapi.NewJSONReader(body),
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error response from search: %s", res.String())
	}

	var searchRes SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}

	return &searchRes, nil
}

// BulkDocument represents a document for bulk indexing
type BulkDocument struct {
	ID       string
	Document interface{}
}

// SearchResponse represents Elasticsearch search response
type SearchResponse struct {
	Took     int64                  `json:"took"`
	TimedOut bool                   `json:"timed_out"`
	Shards   map[string]interface{} `json:"_shards"`
	Hits     struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []Hit   `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]Aggregation `json:"aggregations,omitempty"`
}

// Hit represents a search hit
type Hit struct {
	Index  string                 `json:"_index"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

// Aggregation represents an aggregation result
type Aggregation struct {
	Buckets []Bucket               `json:"buckets,omitempty"`
	Value   interface{}            `json:"value,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// Bucket represents an aggregation bucket
type Bucket struct {
	Key      interface{} `json:"key"`
	DocCount int64       `json:"doc_count"`
}

// Refresh refreshes the index
func (c *Client) Refresh(ctx context.Context, indexName string) error {
	req := esapi.IndicesRefreshRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("error refreshing index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response refreshing index: %s", res.String())
	}

	return nil
}

// GetClient returns the underlying Elasticsearch client
func (c *Client) GetClient() *elasticsearch.Client {
	return c.es
}

// HealthCheck checks Elasticsearch cluster health
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.es.Cluster.Health()
	if err != nil {
		return fmt.Errorf("elasticsearch health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch cluster unhealthy: %s", res.String())
	}

	return nil
}
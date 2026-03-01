package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/elastic/go-elasticsearch/v9"
)

type elasticsearchRetriever struct {
	client   *elasticsearch.TypedClient
	embedder embedding.Embedder
	config   *ElasticsearchRetrieverConfig
}

func NewElasticsearchRetriever(client *elasticsearch.TypedClient, embedder embedding.Embedder, config *ElasticsearchRetrieverConfig) retriever.Retriever {
	return &elasticsearchRetriever{client: client, embedder: embedder, config: config}
}

func embeddingFloat64To32(vec64 []float64) []float32 {
	vec32 := make([]float32, len(vec64))
	for i, v := range vec64 {
		vec32[i] = float32(v)
	}
	return vec32
}

func (r *elasticsearchRetriever) Retrieve(ctx context.Context, query string) (string, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return "", err
	}
	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return "", errors.New("嵌入结果为空，无法进行向量检索")
	}
	embedding := embeddings[0]

	options := r.config.Build(embeddingFloat64To32(embedding))

	searchResp, err := r.client.Search().Index(r.config.IndexName).
		Request(options).Do(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to search docs by vector in es: %s", err)
	}
	var stringsBuilder strings.Builder
	for i, hit := range searchResp.Hits.Hits {
		fmt.Fprintf(&stringsBuilder, "document%d:\n%s\n", i+1, hit.Source_)
	}
	return stringsBuilder.String(), nil
}

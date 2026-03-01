package redisstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/redis/go-redis/v9"
)

type RedisStackRetriever struct {
	client   *redis.Client
	embedder embedding.Embedder
	config   *RedisStackRetrieverConfig
}

func NewRedisStackRetriever(client *redis.Client, embedder embedding.Embedder, config *RedisStackRetrieverConfig) retriever.Retriever {
	return &RedisStackRetriever{client: client, embedder: embedder, config: config}
}

func embeddingFloat64To32(vec64 []float64) []float32 {
	vec32 := make([]float32, len(vec64))
	for i, v := range vec64 {
		vec32[i] = float32(v)
	}
	return vec32
}

func (r *RedisStackRetriever) Retrieve(ctx context.Context, query string) (string, error) {
	embeddings, err := r.embedder.EmbedStrings(ctx, []string{query})
	if err != nil {
		return "", err
	}
	if len(embeddings) == 0 || len(embeddings[0]) == 0 {
		return "", errors.New("嵌入结果为空，无法进行向量检索")
	}
	embedding := embeddings[0]
	finalQuery := fmt.Sprintf("*=>[KNN $K @%s $Vector AS vector_score]", r.config.VectorFieldName)
	vectorBytes := embeddingFloat64To32(embedding)

	options := r.config.Build(vectorBytes)

	result, err := r.client.FTSearchWithArgs(ctx, r.config.IndexName, finalQuery, options).Result()
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}
	var stringsBuilder strings.Builder
	for i, doc := range result.Docs {
		docJSON, err := json.Marshal(&doc)
		if err != nil {
			return "", fmt.Errorf("failed to marshal document to JSON: %w", err)
		}
		fmt.Fprintf(&stringsBuilder, "document%d:\n%s\n", i+1, docJSON)
	}
	return stringsBuilder.String(), nil
}

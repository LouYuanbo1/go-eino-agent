package redisstack

import (
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"
	"github.com/redis/go-redis/v9"
)

type RedisStackRetrieverConfig struct {
	retriever.RetrieverConfig
	SortBy       []redis.FTSearchSortBy `mapstructure:"sort_by"`
	ReturnFields []redis.FTSearchReturn `mapstructure:"return_fields"`
}

func (r *RedisStackRetrieverConfig) Validate() error {
	if r.IndexName == "" {
		return fmt.Errorf("index name is empty")
	}
	if r.VectorFieldName == "" {
		return fmt.Errorf("vector field name is empty")
	}
	if r.K <= 0 {
		return fmt.Errorf("k is invalid")
	}
	return nil
}

func (r *RedisStackRetrieverConfig) Build(vector []float32) *redis.FTSearchOptions {
	options := &redis.FTSearchOptions{
		Params: map[string]any{
			"K":      r.K,
			"Vector": vector,
		},
		DialectVersion: 2,
	}
	if r.SortBy != nil {
		options.SortBy = r.SortBy
	}
	if r.ReturnFields != nil {
		options.Return = r.ReturnFields
	}
	return options
}

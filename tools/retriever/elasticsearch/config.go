package elasticsearch

import (
	"fmt"

	"github.com/LouYuanbo1/go-eino-agent/tools/retriever"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

type ElasticsearchRetrieverConfig struct {
	retriever.RetrieverConfig
	SortBy        *[]types.SortCombinations `mapstructure:"sort_by"`
	ReturnFields  *types.SourceField        `mapstructure:"return_fields"`
	NumCandidates int                       `json:"num_candidates"`
}

func (r *ElasticsearchRetrieverConfig) Validate() error {
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

func (r *ElasticsearchRetrieverConfig) Build(vector []float32) *search.Request {
	options := &search.Request{}
	options.Knn = []types.KnnSearch{
		{
			Field:         r.VectorFieldName,
			QueryVector:   vector,
			K:             &r.K,
			NumCandidates: &r.NumCandidates,
		},
	}
	if r.SortBy != nil {
		options.Sort = *r.SortBy
	}
	if r.ReturnFields != nil {
		options.Source_ = r.ReturnFields
	}
	return options
}

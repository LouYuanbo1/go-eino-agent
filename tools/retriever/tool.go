package retriever

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

type RetrieverParams struct {
	Query string `json:"query" jsonschema:"description=用于检索的查询语句"`
}

func RetrieverFunc[T Retriever](ctx context.Context, retriever T) func(ctx context.Context, params *RetrieverParams) (string, error) {
	return func(ctx context.Context, params *RetrieverParams) (string, error) {
		return retriever.Retrieve(ctx, params.Query)
	}
}

func NewRetrieverTool[T Retriever](ctx context.Context, retriever T) (tool.InvokableTool, error) {
	retrieverTool, err := utils.InferTool(
		"retriever", // tool name
		`Retriever are used to search documents in index by query; 
		they can be used to find similar documents to the query`, // tool description
		RetrieverFunc(ctx, retriever))
	if err != nil {
		return nil, err
	}
	return retrieverTool, nil
}

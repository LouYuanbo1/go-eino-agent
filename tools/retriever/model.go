package retriever

import "context"

type Retriever interface {
	Retrieve(ctx context.Context, query string) (string, error)
}

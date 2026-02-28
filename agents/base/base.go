package agent

import (
	"context"

	"github.com/cloudwego/eino/adk"
)

type Agent interface {
	Run(ctx context.Context, input *adk.AgentInput, options ...adk.AgentRunOption) *adk.AsyncIterator[*adk.AgentEvent]
	OutputMessage(ctx context.Context, input string, withReasoning bool, options ...adk.AgentRunOption)
}

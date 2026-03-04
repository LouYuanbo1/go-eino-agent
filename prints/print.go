package prints

import (
	"fmt"
	"io"
	"log"

	"github.com/cloudwego/eino/adk"
)

func PrintMessages(iter *adk.AsyncIterator[*adk.AgentEvent], opts ...PrintOption) {
	cfg := &printConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("错误: %v", event.Err)
			break
		}

		if event.Output != nil && event.Output.MessageOutput != nil {
			mo := event.Output.MessageOutput
			if mo == nil {
				break
			}
			printMessage(mo, cfg)
		}
	}
}

func printMessage(msg *adk.MessageVariant, cfg *printConfig) {
	if msg == nil {
		return
	}
	if !msg.IsStreaming {
		outputMessage(msg.Message, cfg)
	} else if cfg.withStreaming {
		streamMessage(msg.MessageStream, cfg)
	}
}

func outputMessage(msg adk.Message, cfg *printConfig) {
	if msg == nil {
		return
	}
	if cfg.withReasoning {
		fmt.Println(msg.ReasoningContent)
	}
	fmt.Println(msg.Content)
}

func streamMessage(s adk.MessageStream, cfg *printConfig) {
	defer s.Close()

	for {
		chunk, err := s.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		if chunk.ReasoningContent != "" && cfg.withReasoning {
			fmt.Print(chunk.ReasoningContent)
			continue
		}
		if chunk.Content == "" && len(chunk.ToolCalls) == 0 {
			continue
		}
		fmt.Print(chunk.Content)
	}
}

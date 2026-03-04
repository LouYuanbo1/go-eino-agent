package prints

type printConfig struct {
	withStreaming bool
	withReasoning bool
}

type PrintOption func(*printConfig)

func WithStreaming(withStreaming bool) PrintOption {
	return func(cfg *printConfig) {
		cfg.withStreaming = withStreaming
	}
}

func WithReasoning(withReasoning bool) PrintOption {
	return func(cfg *printConfig) {
		cfg.withReasoning = withReasoning
	}
}

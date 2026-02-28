package prints

type printConfig struct {
	withReasoning bool
}

type PrintOption func(*printConfig)

func WithReasoning(withReasoning bool) PrintOption {
	return func(cfg *printConfig) {
		cfg.withReasoning = withReasoning
	}
}

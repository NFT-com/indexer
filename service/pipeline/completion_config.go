package pipeline

var DefaultCompletionConfig = CompletionConfig{
	DryRun:      false,
	MaxAttempts: 10,
}

type CompletionConfig struct {
	DryRun      bool
	MaxAttempts uint
}

type CompletionOption func(*CompletionConfig)

func WithCompletionDryRun(enabled bool) CompletionOption {
	return func(cfg *CompletionConfig) {
		cfg.DryRun = enabled
	}
}

func WithCompletionMaxAttempts(attempts uint) CompletionOption {
	return func(cfg *CompletionConfig) {
		cfg.MaxAttempts = attempts
	}
}

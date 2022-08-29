package pipeline

var DefaultCompletionConfig = CompletionConfig{
	DryRun:     false,
	MaxRetries: 10,
}

type CompletionConfig struct {
	DryRun     bool
	MaxRetries uint
}

type CompletionOption func(*CompletionConfig)

func WithCompletionDryRun(enabled bool) CompletionOption {
	return func(cfg *CompletionConfig) {
		cfg.DryRun = enabled
	}
}

func WithCompletionMaxRetries(retries uint) CompletionOption {
	return func(cfg *CompletionConfig) {
		cfg.MaxRetries = retries
	}
}

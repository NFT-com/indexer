package pipeline

var DefaultParsingConfig = ParsingConfig{
	DryRun:           false,
	MaxAttempts:      10,
	SanitizeInterval: 1000,
}

type ParsingConfig struct {
	DryRun           bool
	MaxAttempts      uint
	SanitizeInterval uint
}

type ParsingOption func(*ParsingConfig)

func WithParsingDryRun(enabled bool) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.DryRun = enabled
	}
}

func WithParsingMaxAttempts(attempts uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.MaxAttempts = attempts
	}
}

func WithSanitizeInterval(jobs uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.SanitizeInterval = jobs
	}
}

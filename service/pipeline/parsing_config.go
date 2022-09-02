package pipeline

var DefaultParsingConfig = ParsingConfig{
	DryRun:           false,
	MaxAttempts:      3,
	SplitRatio:       2,
	SanitizeInterval: 10000,
}

type ParsingConfig struct {
	DryRun           bool
	MaxAttempts      uint16
	SplitRatio       uint
	SanitizeInterval uint
}

type ParsingOption func(*ParsingConfig)

func WithParsingDryRun(enabled bool) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.DryRun = enabled
	}
}

func WithParsingMaxAttempts(attempts uint16) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.MaxAttempts = attempts
	}
}

func WithSplitRatio(ratio uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.SplitRatio = ratio
	}
}

func WithSanitizeInterval(jobs uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.SanitizeInterval = jobs
	}
}

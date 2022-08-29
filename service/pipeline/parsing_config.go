package pipeline

var DefaultParsingConfig = ParsingConfig{
	DryRun:           false,
	MaxRetries:       10,
	SanitizeInterval: 1000,
}

type ParsingConfig struct {
	DryRun           bool
	MaxRetries       uint
	SanitizeInterval uint
}

type ParsingOption func(*ParsingConfig)

func WithParsingDryRun(enabled bool) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.DryRun = enabled
	}
}

func WithParsingMaxRetries(retries uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.MaxRetries = retries
	}
}

func WithSanitizeInterval(numJobs uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.SanitizeInterval = numJobs
	}
}

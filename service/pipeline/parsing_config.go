package pipeline

var DefaultParsingConfig = ParsingConfig{
	DryRun:           false,
	SanitizeInterval: 1000,
}

type ParsingConfig struct {
	DryRun           bool
	SanitizeInterval uint
}

type ParsingOption func(*ParsingConfig)

func WithDryRun(enabled bool) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.DryRun = enabled
	}
}

func WithSanitizeInterval(numJobs uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.SanitizeInterval = numJobs
	}
}

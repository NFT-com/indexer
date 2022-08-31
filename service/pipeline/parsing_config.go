package pipeline

var DefaultParsingConfig = ParsingConfig{
	DryRun:           false,
	MaxAttempts:      10,
	MaxAddresses:     100,
	MaxHeights:       10,
	SplitRatio:       2,
	SanitizeInterval: 10000,
}

type ParsingConfig struct {
	DryRun           bool
	MaxAttempts      uint16
	MaxAddresses     uint
	MaxHeights       uint
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

func WithParsingMaxAddresses(addresses uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.MaxAddresses = addresses
	}
}

func WithParsingMaxHeights(heights uint) ParsingOption {
	return func(cfg *ParsingConfig) {
		cfg.MaxHeights = heights
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

package pipeline

var DefaultAdditionConfig = AdditionConfig{
	DryRun:      false,
	MaxAttempts: 10,
}

type AdditionConfig struct {
	DryRun      bool
	MaxAttempts uint
}

type AdditionOption func(*AdditionConfig)

func WithAdditionDryRun(enabled bool) AdditionOption {
	return func(cfg *AdditionConfig) {
		cfg.DryRun = enabled
	}
}

func WithAdditionMaxAttempts(attempts uint) AdditionOption {
	return func(cfg *AdditionConfig) {
		cfg.MaxAttempts = attempts
	}
}

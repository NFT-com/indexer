package pipeline

var DefaultAdditionConfig = AdditionConfig{
	DryRun:     false,
	MaxRetries: 10,
}

type AdditionConfig struct {
	DryRun     bool
	MaxRetries uint
}

type AdditionOption func(*AdditionConfig)

func WithAdditionDryRun(enabled bool) AdditionOption {
	return func(cfg *AdditionConfig) {
		cfg.DryRun = enabled
	}
}

func WithAdditionMaxRetries(retries uint) AdditionOption {
	return func(cfg *AdditionConfig) {
		cfg.MaxRetries = retries
	}
}

package web2

var MetadataDefaultConfig = MetadataConfig{
	DisableValidation: false,
}

type MetadataConfig struct {
	DisableValidation bool
}

type MetadataOption func(*MetadataConfig)

func WithDisableValidation(disabled bool) MetadataOption {
	return func(cfg *MetadataConfig) {
		cfg.DisableValidation = disabled
	}
}

package config

// Core provider config

type Config struct {
	PollInterval string           `json:"pollInterval,omitempty" yaml:"pollInterval,omitempty"`
	Providers    []ProviderConfig `json:"providers,omitempty" yaml:"providers,omitempty"`
}

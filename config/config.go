package config

// Core provider config

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		PollInterval: "5s",
		Providers: []ProviderConfig{
			{
				Name: "Traefik Provider",
				Connection: ConnectionConfig{
					Host:    []string{"localhost"},
					Port:    8080,
					Path:    "/api/rawdata",
					Timeout: "5s",
					Headers: map[string]string{
						"Host": "dashboard.traefik.localhost",
					},
					MTLS: nil,
				},
			},
		},
	}
}

type Config struct {
	PollInterval string           `json:"pollInterval,omitempty" yaml:"pollInterval,omitempty"`
	Providers    []ProviderConfig `json:"providers,omitempty" yaml:"providers,omitempty"`
}

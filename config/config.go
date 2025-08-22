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
				HTTP: &HTTPSection{
					Discover: true,
					Routers: &RoutersConfig{
						Discover:         true,
						DiscoverPriority: false,
					},
					Services: &ServicesConfig{
						Discover: true,
					},
					Middlewares: &MiddlewaresConfig{
						Discover: true,
					},
					ServerTransports: ServerTransportsConfig{
						Discover: true,
					},
				},
				TCP: &TCPSection{
					Discover: true,
					Routers: &RoutersConfig{
						Discover:         true,
						DiscoverPriority: false,
					},
					Middlewares: &MiddlewaresConfig{
						Discover: true,
					},
					Services: &ServicesConfig{
						Discover: true,
					},
				},
				UDP: &UDPSection{
					Discover: true,
					Routers: &UDPRoutersConfig{
						Discover: true,
					},
					Services: &UDPServicesConfig{
						Discover: true,
					},
				},
				TLS: &TLSSection{
					Discover: true,
				},
			},
		},
	}
}

type Config struct {
	PollInterval string           `json:"pollInterval,omitempty" yaml:"pollInterval,omitempty"`
	Providers    []ProviderConfig `json:"providers,omitempty" yaml:"providers,omitempty"`
}

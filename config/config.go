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
					Host:    "localhost",
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
						Filters: RouterFilters{
							Name:        "",
							Entrypoints: []string{},
							Rule:        "",
							Service:     "",
						},
					},
					Services: &ServicesConfig{
						Discover: true,
					},
					Middlewares: &MiddlewaresConfig{
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

type ProviderConfig struct {
	Name       string           `json:"name,omitempty" yaml:"name,omitempty"`
	Connection ConnectionConfig `json:"connection,omitempty" yaml:"connection,omitempty"`
	HTTP       *HTTPSection     `json:"http,omitempty" yaml:"http,omitempty"`
	TCP        *TCPSection      `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP        *UDPSection      `json:"udp,omitempty" yaml:"udp,omitempty"`
	TLS        *TLSSection      `json:"tls,omitempty" yaml:"tls,omitempty"`
	Tunnels    []TunnelConfig   `json:"tunnels,omitempty" yaml:"tunnels,omitempty"`
}

type ConnectionConfig struct {
	Host    string            `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int               `json:"port,omitempty" yaml:"port,omitempty"`
	Path    string            `json:"path,omitempty" yaml:"path,omitempty"`
	Timeout string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	MTLS    *MTLSConfig       `json:"mTLS,omitempty" yaml:"mTLS,omitempty"`
}

type MTLSConfig struct {
	CAFile   string `json:"caFile,omitempty" yaml:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" yaml:"keyFile,omitempty"`
}

type TunnelConfig struct {
	Name      string      `json:"name,omitempty" yaml:"name,omitempty"`
	Addresses []string    `json:"connection,omitempty" yaml:"connection,omitempty"`
	MTLS      *MTLSConfig `json:"mTLS,omitempty" yaml:"mTLS,omitempty"`
}

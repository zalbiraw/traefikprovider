// Package config defines the provider configuration model and parsing helpers.
package config

//revive:disable:tagliatelle

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
				},
				Filter: ProviderFilter{
					Provider: "file",
				},
			},
		},
	}
}

// Config is the root configuration for the provider plugin.
// It controls the polling interval and the list of upstream providers
// to fetch and filter dynamic configuration from.
type Config struct {
	PollInterval string           `json:"pollInterval,omitempty" yaml:"pollInterval,omitempty"`
	Providers    []ProviderConfig `json:"providers,omitempty" yaml:"providers,omitempty"`
}

// ProviderConfig defines a single upstream provider to poll and filter.
type ProviderConfig struct {
	Name       string           `json:"name,omitempty" yaml:"name,omitempty"`
	Filter     ProviderFilter   `json:"filter,omitempty" yaml:"filter,omitempty"`
	Connection ConnectionConfig `json:"connection,omitempty" yaml:"connection,omitempty"`
	HTTP       *HTTPSection     `json:"http,omitempty" yaml:"http,omitempty"`
	TCP        *TCPSection      `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP        *UDPSection      `json:"udp,omitempty" yaml:"udp,omitempty"`
	TLS        *TLSSection      `json:"tls,omitempty" yaml:"tls,omitempty"`
	Tunnels    []TunnelConfig   `json:"tunnels,omitempty" yaml:"tunnels,omitempty"`
}

// ProviderFilter narrows discovered objects to a specific provider name.
type ProviderFilter struct {
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}

// ConnectionConfig configures how to connect to the upstream provider API.
type ConnectionConfig struct {
	Host    string            `json:"host,omitempty" yaml:"host,omitempty"`
	Port    int               `json:"port,omitempty" yaml:"port,omitempty"`
	Path    string            `json:"path,omitempty" yaml:"path,omitempty"`
	Timeout string            `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	MTLS    *MTLSConfig       `json:"mTLS,omitempty" yaml:"mTLS,omitempty"` //nolint:tagliatelle
}

// MTLSConfig holds mutual TLS file paths for establishing mTLS connections.
type MTLSConfig struct {
	CAFile   string `json:"caFile,omitempty" yaml:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty" yaml:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" yaml:"keyFile,omitempty"`
}

// TunnelConfig defines a list of addresses grouped under a name, optionally with mTLS.
type TunnelConfig struct {
	Name      string      `json:"name,omitempty" yaml:"name,omitempty"`
	Addresses []string    `json:"connection,omitempty" yaml:"connection,omitempty"`
	MTLS      *MTLSConfig `json:"mTLS,omitempty" yaml:"mTLS,omitempty"` //nolint:tagliatelle
}

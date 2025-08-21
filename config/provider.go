package config

type ProviderConfig struct {
	Name        string             `json:"name,omitempty" yaml:"name,omitempty"`
	Connection  ConnectionConfig   `json:"connection,omitempty" yaml:"connection,omitempty"`
	HTTP        *HTTPSection       `json:"http,omitempty" yaml:"http,omitempty"`
	TCP         *TCPSection        `json:"tcp,omitempty" yaml:"tcp,omitempty"`
	UDP         *UDPSection        `json:"udp,omitempty" yaml:"udp,omitempty"`
	TLS         *TLSSection        `json:"tls,omitempty" yaml:"tls,omitempty"`
}

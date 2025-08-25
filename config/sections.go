package config

// HTTPSection controls discovery of HTTP routers, middlewares, and services.
type HTTPSection struct {
	Discover    bool               `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers     *RoutersConfig     `json:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares *MiddlewaresConfig `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services    *ServicesConfig    `json:"services,omitempty" yaml:"services,omitempty"`
}

// TCPSection controls discovery of TCP routers, middlewares, and services.
type TCPSection struct {
	Discover    bool               `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers     *RoutersConfig     `json:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares *MiddlewaresConfig `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services    *ServicesConfig    `json:"services,omitempty" yaml:"services,omitempty"`
}

// UDPSection controls discovery of UDP routers and services.
type UDPSection struct {
	Discover bool               `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers  *UDPRoutersConfig  `json:"routers,omitempty" yaml:"routers,omitempty"`
	Services *UDPServicesConfig `json:"services,omitempty" yaml:"services,omitempty"`
}

// TLSSection controls discovery of TLS options, stores, and certificates.
type TLSSection struct {
	Discover bool `json:"discover,omitempty" yaml:"discover,omitempty"`
}

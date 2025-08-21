package config

type HTTPSection struct {
	Discover         bool                    `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers          *RoutersConfig          `json:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares      *MiddlewaresConfig      `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services         *ServicesConfig         `json:"services,omitempty" yaml:"services,omitempty"`
	ServerTransports ServerTransportsConfig `json:"serverTransports,omitempty" yaml:"serverTransports,omitempty"`
}

type TCPSection struct {
	Discover    bool               `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers     *RoutersConfig     `json:"routers,omitempty" yaml:"routers,omitempty"`
	Middlewares *MiddlewaresConfig `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
	Services    *ServicesConfig    `json:"services,omitempty" yaml:"services,omitempty"`
}

type UDPSection struct {
	Discover bool               `json:"discover,omitempty" yaml:"discover,omitempty"`
	Routers  *UDPRoutersConfig  `json:"routers,omitempty" yaml:"routers,omitempty"`
	Services *UDPServicesConfig `json:"services,omitempty" yaml:"services,omitempty"`
}

type TLSSection struct {
	Discover bool `json:"discover,omitempty" yaml:"discover,omitempty"`
}

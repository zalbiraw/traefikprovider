package config

// UDPRoutersConfig holds discovery, filtering, and overrides for UDP routers.
type UDPRoutersConfig struct {
	Discover    bool            `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter      UDPRouterFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	Overrides   UDPOverrides    `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes []interface{}   `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

// UDPRouterFilter filters UDP routers by name, provider, entrypoints, and service.
type UDPRouterFilter struct {
	Name        string   `json:"name,omitempty" yaml:"name,omitempty"`
	Provider    string   `json:"provider,omitempty" yaml:"provider,omitempty"`
	Entrypoints []string `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Service     string   `json:"service,omitempty" yaml:"service,omitempty"`
}

// UDPOverrides defines overrides applied to filtered UDP routers.
type UDPOverrides struct {
	Name        string               `json:"name,omitempty" yaml:"name,omitempty"`
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
}

// UDPServicesConfig holds discovery, filtering, and overrides for UDP services.
type UDPServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter        ServiceFilter    `json:"filter,omitempty" yaml:"filter,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

package config

type UDPRoutersConfig struct {
	Discover    bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters     UDPRouterFilters `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides   UDPOverrides     `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes []interface{}    `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

type UDPRouterFilters struct {
	Name         string   `json:"name,omitempty" yaml:"name,omitempty"`
	NameRegex    string   `json:"nameRegex,omitempty" yaml:"nameRegex,omitempty"`
	Entrypoints  []string `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Service      string   `json:"service,omitempty" yaml:"service,omitempty"`
	ServiceRegex string   `json:"serviceRegex,omitempty" yaml:"serviceRegex,omitempty"`
}

type UDPOverrides struct {
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
}

type UDPServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filters       ServiceFilters   `json:"filters,omitempty" yaml:"filters,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

package config

// UDPRoutersConfig holds discovery, matching, and overrides for UDP routers.
type UDPRoutersConfig struct {
	Discover    bool          `json:"discover,omitempty" yaml:"discover,omitempty"`
	Matcher     string        `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	Overrides   UDPOverrides  `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraRoutes []interface{} `json:"extraRoutes,omitempty" yaml:"extraRoutes,omitempty"`
}

// UDPOverrides defines overrides applied to matched UDP routers.
type UDPOverrides struct {
	Name        string               `json:"name,omitempty" yaml:"name,omitempty"`
	Entrypoints []OverrideEntrypoint `json:"entrypoints,omitempty" yaml:"entrypoints,omitempty"`
	Services    []OverrideService    `json:"services,omitempty" yaml:"services,omitempty"`
}

// UDPServicesConfig holds discovery, matching, and overrides for UDP services.
type UDPServicesConfig struct {
	Discover      bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Matcher       string           `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	Overrides     ServiceOverrides `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	ExtraServices []interface{}    `json:"extraServices,omitempty" yaml:"extraServices,omitempty"`
}

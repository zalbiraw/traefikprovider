package config

// MiddlewaresConfig holds discovery and filter settings for middlewares.
type MiddlewaresConfig struct {
	Discover         bool             `json:"discover,omitempty" yaml:"discover,omitempty"`
	Filter           MiddlewareFilter `json:"filter,omitempty" yaml:"filter,omitempty"`
	ExtraMiddlewares []interface{}    `json:"extraMiddlewares,omitempty" yaml:"extraMiddlewares,omitempty"`
}

// MiddlewareFilter filters middlewares by name and provider.
type MiddlewareFilter struct {
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Provider string `json:"provider,omitempty" yaml:"provider,omitempty"`
}

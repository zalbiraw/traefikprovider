package config

// MiddlewaresConfig holds discovery and matcher settings for middlewares.
type MiddlewaresConfig struct {
	Discover         bool          `json:"discover,omitempty" yaml:"discover,omitempty"`
	Matcher          string        `json:"matcher,omitempty" yaml:"matcher,omitempty"`
	ExtraMiddlewares []interface{} `json:"extraMiddlewares,omitempty" yaml:"extraMiddlewares,omitempty"`
}

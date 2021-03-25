package config

var (
	ConfigFileName = "config.yaml"
)

type Import struct {
    kind string `json:"kind"`
}

type Service struct {
	Import []Import `json:"import"`
}

type Config struct {
	Service Service `json:"service"`
}

func New() (*Config, error) {
	c := &Config{}
	if err := c.Load(); err != nil {
		return c, err
	}

	return c, nil
}

func (c *Config) Load() error {
	// Use Viper to Load config?
	return nil
}
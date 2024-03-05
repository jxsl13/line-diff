package config

type Config struct {
	Sorted bool `koanf:"sorted" short:"s" description:"sort the results"`
}

func (c *Config) Validate() error {
	return nil
}

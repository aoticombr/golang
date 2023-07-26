package jsonconfig

type Config struct {
	Boots []Boot `json:"boots"`
	Apis  []Api  `json:"apis"`
	Path  string `json:"path"`
}

func (c *Config) GetPath() string {
	return c.Path
}
func (c *Config) GetBoots() []Boot {
	return c.Boots
}
func (c *Config) GetApis() []Api {
	return c.Apis
}
func (c *Config) GetBootByName(name string) *Boot {
	for _, boot := range c.Boots {
		if boot.Name == name {
			return &boot
		}
	}
	return nil
}
func (c *Config) GetApiByName(name string) *Api {
	for _, api := range c.Apis {
		if api.Name == name {
			return &api
		}
	}
	return nil
}

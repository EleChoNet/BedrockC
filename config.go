package bedrockc
import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/pkg/errors"

  )
const ConfigPath = "./config.yaml"


type Config struct {
	ConfigFile string
	values map[string]string
}
func NewConfig() (*Config, error) {
	c := &Config{}
	c.values = make(map[string]string)
	return c, nil
}
func (c *Config) Load() error {
	data, err := ioutil.ReadFile(c.ConfigFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	err = yaml.Unmarshal(data, &c.values)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}
	return nil
}
func (c *Config) Require(key string,defaultValue string) string{
	if _, ok := c.values[key]; !ok {
		c.values[key] = defaultValue
	}
	return c.values[key]
}


func (c Config) Save() error {
	data, err := yaml.Marshal(c.values)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}
	return ioutil.WriteFile(c.ConfigFile, data, 0644)
}

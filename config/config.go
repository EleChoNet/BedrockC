package config

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)
const ConfigPath = "./config.yaml"


type Config struct {
	ConfigFile string
	values map[string]interface{}
}
func NewConfig() (*Config, error) {
	c := &Config{}
	c.values = make(map[string]interface{})
	c.ConfigFile = ConfigPath
	//如果配置文件不存在，则创建配置文件
	if _, err := os.Stat(c.ConfigFile); os.IsNotExist(err) {
		err := c.Save()
		if err != nil {
			return nil, errors.Wrap(err, "failed to save config")
		}
	}
	
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
func (c *Config) Require(key string,defaultValue interface{}) interface{}{
	if _, ok := c.values[key]; !ok {
		c.values[key] = defaultValue
	}
	return c.values[key]
}
func (c *Config) Set(key string,value interface{}){
	c.values[key] = value
}


func (c Config) Save() error {
	data, err := yaml.Marshal(c.values)
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}
	return ioutil.WriteFile(c.ConfigFile, data, 0644)
}

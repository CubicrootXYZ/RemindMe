package configuration

import (
	"github.com/jinzhu/configor"
)

// Load loads the given configuration files
func Load(files []string) (*Config, error) {
	config := &Config{}

	err := configor.New(&configor.Config{
		Environment:          "production",
		ENVPrefix:            "REMINDME",
		ErrorOnUnmatchedKeys: false,
	}).Load(config, files...)

	return config, err
}

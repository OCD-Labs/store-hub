package util

import "github.com/spf13/viper"

// A Configs defines the expected config values.
type Configs struct {
	Env                string   `mapstructure:"ENVIRONMENT"`
	ServerAddr         string   `mapstructure:"SERVER_ADDRESS"`
	CorsTrustedOrigins []string `mapstructure:"CORS_TRUSTED_ORIGINS"`
	DBDriver           string   `mapstructure:"DB_DRIVER"`
	DBSource           string   `mapstructure:"DB_SOURCE"`
}

// ParseConfigs parses the configuration files.
func ParseConfigs(path string) (config Configs, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("secrets")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
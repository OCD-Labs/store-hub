package util

import "github.com/spf13/viper"

// A Configs defines the expected config values.
type Configs struct {
	Env                 string   `mapstructure:"ENVIRONMENT"`
	ServerAddr          string   `mapstructure:"SERVER_ADDRESS"`
	CorsTrustedOrigins  []string `mapstructure:"CORS_TRUSTED_ORIGINS"`
	DBDriver            string   `mapstructure:"DB_DRIVER"`
	DBSource            string   `mapstructure:"DB_SOURCE"`
	EmailSenderAddress  string   `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderName     string   `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderPassword string   `mapstructure:"EMAIL_SENDER_PASSWORD"`
	TokenSymmetricKey   string   `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RedisAddress        string   `mapstructure:"REDIS_ADDRESS"`
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
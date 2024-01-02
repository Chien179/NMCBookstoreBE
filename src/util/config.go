package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or enviroment variables.
type Config struct {
	DBDriver                string        `mapstructure:"DB_DRIVER"`
	DBSource                string        `mapstructure:"DB_SOURCE"`
	DBSourceTest            string        `mapstructure:"DB_SOURCE_TEST"`
	ServerAddress           string        `mapstructure:"SERVER_ADDRESS"`
	GRPCAddress             string        `mapstructure:"GRPC_ADDRESS"`
	TokenSymmetrictKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration     time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration    time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MigrationURL            string        `mapstructure:"MIGRATION_URL"`
	Environment             string        `mapstructure:"ENVIRONMENT"`
	EmailSenderName         string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress      string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword     string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	RedisAddress            string        `mapstructure:"REDIS_ADDRESS"`
	CloudinaryCloudName     string        `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CloudinaryAPIKey        string        `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryAPISecret     string        `mapstructure:"CLOUDINARY_API_SECRET"`
	RedisUsername           string        `mapstructure:"REDIS_USERNAME"`
	RedisPassword           string        `mapstructure:"REDIS_PASSWORD"`
	GoogleOauthClientID     string        `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOauthClientSecret string        `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	ClientHost              string        `mapstructure:"CLIENT_HOST"`
	ElasticUsername         string        `mapstructure:"ELASTIC_USERNAME"`
	ElasticAddress          string        `mapstructure:"ELASTIC_ADDRESS"`
	ElasticPassword         string        `mapstructure:"ELASTIC_PASSWORD"`
}

// LoadConfig reads configuration from file or enviroment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

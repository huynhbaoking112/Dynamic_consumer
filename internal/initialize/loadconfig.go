package initialize

import (
	"event_service/global"
	"event_service/internal/common"
	"event_service/pkg/setting"
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() {
	fmt.Println("Loading configuration...")

	viper := viper.New()
	viper.AddConfigPath("configs")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("EVENT_SERVICE")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		fmt.Println("Using default configuration...")
		loadDefaultConfig()
		return
	}

	var config setting.Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshalling config: %v\n", err)
		panic(fmt.Errorf("%w: %v", common.ErrConfigLoad, err))
	}

	if err := validateConfig(&config); err != nil {
		panic(fmt.Errorf("%w: %v", common.ErrConfigValidation, err))
	}

	global.Config = &config
	fmt.Printf("Configuration loaded successfully from: %s\n", viper.ConfigFileUsed())
}

func loadDefaultConfig() {
	config := &setting.Config{
		Server: setting.Server{
			Host: "localhost",
			Port: 8081,
		},
		MongoDB: setting.MongoDB{
			Host:                  "localhost",
			Port:                  27017,
			Database:              "notification",
			ActivityLogCollection: "activity_logs",
			User:                  "admin",
			Password:              "password",
			ConnectionString:      "mongodb://admin:password@localhost:27017/notification",
		},
		RabbitMQ: setting.RabbitMQ{
			Host:                  "localhost",
			Port:                  5672,
			User:                  "guest",
			Password:              "guest",
			IAMExchange:           "iam_events_topic",
			ActivityLogQueue:      "iam_activity_log_queue",
			ActivityLogBindingKey: "#.log",
			RetryAttempts:         3,
			RetryDelaySeconds:     5,
		},
	}

	global.Config = config
	fmt.Println("Default configuration loaded")
}

func validateConfig(config *setting.Config) error {
	if config.MongoDB.Database == "" {
		return fmt.Errorf("mongodb database name is required")
	}

	if config.MongoDB.ActivityLogCollection == "" {
		return fmt.Errorf("mongodb activity log collection name is required")
	}

	if config.RabbitMQ.IAMExchange == "" {
		return fmt.Errorf("rabbitmq iam exchange name is required")
	}

	if config.RabbitMQ.ActivityLogQueue == "" {
		return fmt.Errorf("rabbitmq activity log queue name is required")
	}

	if config.RabbitMQ.ActivityLogBindingKey == "" {
		return fmt.Errorf("rabbitmq binding key is required")
	}

	return nil
}

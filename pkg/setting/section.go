package setting

// MongoDB configuration
type MongoDB struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	Database              string `mapstructure:"database"`
	ActivityLogCollection string `mapstructure:"activity_log_collection"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	ConnectionString      string `mapstructure:"connection_string"`
}

// RabbitMQ configuration
type RabbitMQ struct {
	Host                  string `mapstructure:"host"`
	Port                  int    `mapstructure:"port"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"password"`
	IAMExchange           string `mapstructure:"iam_exchange"`
	AppStoreExchange      string `mapstructure:"app_store_exchange"`
	ActivityLogQueue      string `mapstructure:"activity_log_queue"`
	ActivityLogBindingKey string `mapstructure:"activity_log_binding_key"`
	RetryAttempts         int    `mapstructure:"retry_attempts"`
	RetryDelaySeconds     int    `mapstructure:"retry_delay_seconds"`
}

// Server configuration
type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Main configuration struct
type Config struct {
	Server   Server   `mapstructure:"server"`
	MongoDB  MongoDB  `mapstructure:"mongodb"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

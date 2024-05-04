package apiserver

type DatabaseURL struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"db_name"`
	FullName string `yaml:"fullname"`
}

type RedisURL struct {
	LocalHost string `yaml:"localhost"`
	Docker    string `yaml:"docker"`
}

type Config struct {
	BindAddr      string      `yaml:"bind_addr"`
	LogLevel      string      `yaml:"log_level"`
	DatabaseURL   DatabaseURL `yaml:"database_url"`
	SessionKey    string      `yaml:"session_key"`
	RedisURL      RedisURL    `yaml:"redis_url"`
	LocalHostMode bool        `yaml:"localhost_mode"`
}

func NewConfig() Config {
	return Config{
		// BindAddr: ":8080",
		// LogLevel: "debug",
	}
}

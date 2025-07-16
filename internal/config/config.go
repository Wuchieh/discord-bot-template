package config

var (
	cfg *Config
)

type DB struct {
	File     string `yaml:"file"`
	LogLevel string `yaml:"log_level"`
}

type Config struct {
	Token string `yaml:"bot_token"`
	DB    DB     `yaml:"db"`
}

func GetDefault() Config {
	return Config{
		Token: "",
		DB: DB{
			File:     "database.db",
			LogLevel: "warn",
		},
	}
}

func Get() Config {
	if cfg == nil {
		cfg = new(Config)
	}

	return *cfg
}

func Set(c Config) {
	cfg = &c
}

package config

var (
	cfg *Config
)

type Config struct {
	Token  string `yaml:"bot_token"`
	DBFile string `yaml:"db_file"`
}

func GetDefault() Config {
	return Config{
		Token:  "",
		DBFile: "database.db",
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

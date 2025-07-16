package config

import (
	"github.com/Wuchieh/discord-bot-template/internal/bot"
	"github.com/Wuchieh/discord-bot-template/internal/bot/handler/reaction_role"
	"github.com/Wuchieh/discord-bot-template/internal/database"
)

var (
	cfg *Config
)

type Config struct {
	Token string `yaml:"bot_token"`

	DB database.Config `yaml:"db"`

	ReactionRole reaction_role.Config `yaml:"reaction_role"`
}

func GetDefault() Config {
	return Config{
		Token: bot.DefaultToken,
		DB: database.Config{
			File:     "database.db",
			LogLevel: "warn",
		},
		ReactionRole: reaction_role.Config{
			GuildID: []string{
				"GuildID1",
				"GuildID2",
				"GuildID3",
			},
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

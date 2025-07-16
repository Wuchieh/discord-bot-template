package model

import "gorm.io/gorm"

type ReactionRole struct {
	gorm.Model
	GuildID   string `gorm:"column:guild_id;uniqueIndex:uni_gme"`
	MessageID string `gorm:"column:message_id;uniqueIndex:uni_gme"`
	Emoji     string `gorm:"column:emoji;uniqueIndex:uni_gme;comment:表情符號的ApiName"`
	RoleID    string `gorm:"column:role_id"`
}

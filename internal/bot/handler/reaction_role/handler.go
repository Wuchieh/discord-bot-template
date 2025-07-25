package reaction_role

import (
	"fmt"
	"github.com/Wuchieh/discord-bot-template/internal/bot/handler"
	"github.com/Wuchieh/discord-bot-template/internal/database"
	"github.com/Wuchieh/discord-bot-template/internal/model"
	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
	"strings"
)

const (
	CommandName = "reaction_role"
)

type Config struct {
	Enabled bool     `yaml:"enabled"`
	GuildID []string `yaml:"guild_id"`
}

var (
	defaultMemberPermissions = int64(discordgo.PermissionManageGuild)
	registeredCommands       = make(map[string]*discordgo.ApplicationCommand)
	command                  = &discordgo.ApplicationCommand{
		Name:                     CommandName,
		DefaultMemberPermissions: &defaultMemberPermissions,
		Description:              "使用反應添加身分組",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message_id",
				Description: "訊息ID",
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
				Required: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "emoji",
				Description: "表情符號",
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
				Required: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "身分組",
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
				Required:  true,
				MaxLength: 1,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "force",
				Description: "強制新增(針對特殊情況設置)",
				ChannelTypes: []discordgo.ChannelType{
					discordgo.ChannelTypeGuildText,
				},
				Required:  false,
				MaxLength: 1,
			},
		},
	}
)

func ptr(s string) *string {
	return &s
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	appData := i.ApplicationCommandData()
	if appData.Name != CommandName {
		return
	}

	// 請稍等...
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}); err != nil {
		fmt.Println("InteractionRespondError:", err)
		return
	}

	guildID := i.GuildID
	messageID := appData.GetOption("message_id").StringValue()

	force := false
	forceOption := appData.GetOption("force")
	if forceOption != nil && forceOption.Type == discordgo.ApplicationCommandOptionBoolean {
		force = forceOption.Value == true
	}

	emoji := appData.GetOption("emoji").StringValue()
	if strings.TrimSpace(emoji) == "" {
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptr("尚未選擇表情符號"),
		}); err != nil {
			fmt.Println("InteractionResponseEditError:", err)
		}
		return
	}

	var emojiApiName string
	emoji = strings.TrimPrefix(emoji, "<")
	emoji = strings.TrimSuffix(emoji, ">")
	emojiSlice := strings.Split(emoji, ":")

	switch len(emojiSlice) {
	case 1:
		emojiApiName = emoji
	case 2:
		emojiApiName = strings.Join(emojiSlice[0:2], ":")
	case 3:
		emojiApiName = strings.Join(emojiSlice[1:3], ":")
	default:
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptr("表情符號格式錯誤"),
		}); err != nil {
			fmt.Println("InteractionResponseEditError:", err)
		}
		return
	}

	role := appData.GetOption("role").RoleValue(s, i.GuildID)
	if role == nil {
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptr("無法取得身分組"),
		}); err != nil {
			fmt.Println("InteractionResponseEditError:", err)
		}
		return
	}

	// 取得訊息
	message, err := s.ChannelMessage(i.ChannelID, messageID)
	if err != nil {
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptr("無法去得指定訊息"),
		}); err != nil {
			fmt.Println("InteractionResponseEditError:", err)
		}
		return
	}

	// 檢查是否已存在相同的emoji
	for _, reaction := range message.Reactions {
		if reaction.Emoji.APIName() == emojiApiName {
			if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: ptr("已存在相同的 emoji"),
			}); err != nil {
				fmt.Println("InteractionResponseEditError:", err)
			}
			return
		}
	}

	// 建立紀錄
	db := database.GetDB()
	if err := db.Transaction(func(tx *gorm.DB) error {
		if force {
			var record model.ReactionRole
			tx.Where(
				model.ReactionRole{
					GuildID:   guildID,
					MessageID: messageID,
					Emoji:     emojiApiName,
				},
			).First(&record)

			record.GuildID = guildID
			record.MessageID = messageID
			record.Emoji = emojiApiName
			record.RoleID = role.ID

			if err := tx.Save(&record).Error; err != nil {
				return fmt.Errorf("創建紀錄失敗: %v", err)
			}
		} else {
			// 創建紀錄
			m := model.ReactionRole{
				GuildID:   guildID,
				MessageID: messageID,
				Emoji:     emojiApiName,
				RoleID:    role.ID,
			}

			if err := tx.Save(&m).Error; err != nil {
				return fmt.Errorf("創建紀錄失敗: %v", err)
			}
		}

		// 在該訊息設置Emoji
		if err := s.MessageReactionAdd(i.ChannelID, messageID, emojiApiName); err != nil {
			return fmt.Errorf("訊息設置Emoji失敗 %v", err)
		}

		return nil
	}); err != nil {
		if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: ptr(fmt.Sprintf("設置失敗: %v", err)),
		}); err != nil {
			fmt.Println("InteractionResponseEditError:", err)
		}
		return
	}

	// 配置完畢
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: ptr("success"),
	}); err != nil {
		fmt.Println("InteractionResponseEditError:", err)
	}
}

// 用戶點下 emoji
func addReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.Member.User.Bot {
		return
	}

	guildID := r.GuildID
	messageID := r.MessageID
	emoji := r.Emoji.APIName()

	var record model.ReactionRole
	db := database.GetDB()
	db.Where(
		model.ReactionRole{
			GuildID:   guildID,
			MessageID: messageID,
			Emoji:     emoji,
		},
	).First(&record)

	if record.ID == 0 {
		return
	}

	err := s.GuildMemberRoleAdd(guildID, r.UserID, record.RoleID)
	if err != nil {
		fmt.Println("設定身分組失敗:", err)
	}
}

// 用戶移除 emoji
func removeReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	guildID := r.GuildID
	messageID := r.MessageID
	emoji := r.Emoji.APIName()

	var record model.ReactionRole
	db := database.GetDB()
	db.Where(
		model.ReactionRole{
			GuildID:   guildID,
			MessageID: messageID,
			Emoji:     emoji,
		},
	).First(&record)

	if record.ID == 0 {
		return
	}

	err := s.GuildMemberRoleRemove(guildID, r.UserID, record.RoleID)
	if err != nil {
		fmt.Println("移除身分組失敗:", err)
	}
}

// 管理員移除所有 emoji
func removeAllReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemoveAll) {
	guildID := r.GuildID
	messageID := r.MessageID
	emoji := r.Emoji.APIName()
	db := database.GetDB()
	if err := db.
		Unscoped().
		Where(
			model.ReactionRole{
				GuildID:   guildID,
				MessageID: messageID,
				Emoji:     emoji,
			},
		).
		Delete(&model.ReactionRole{}).
		Error; err != nil {
		fmt.Println("資料庫資料刪除失敗:", err)
	}
}

func Setup(cfg Config) {
	if !cfg.Enabled {
		return
	}

	fmt.Printf("加載模組: reaction_role, 指令: /%s\n", CommandName)

	handler.OnOpened(func(s *discordgo.Session) {
		for _, guildID := range cfg.GuildID {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, command)
			if err != nil {
				fmt.Printf("GuildID: %s, %s 指令註冊失敗\n", guildID, CommandName)
			} else {
				registeredCommands[guildID] = cmd
			}
		}
	})

	handler.OnBeforeClose(func(s *discordgo.Session) {
		for guildID, cmd := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, guildID, cmd.ID)
			if err != nil {
				fmt.Printf("GuildID: %s, %s 指令移除失敗\n", guildID, CommandName)
			}
		}
	})

	handler.AddHandler(
		commandHandler,
		addReactionHandler,
		removeReactionHandler,
		removeAllReactionHandler,
	)
}

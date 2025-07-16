package bootstrap

import (
	"errors"
	"fmt"
	"github.com/Wuchieh/discord-bot-template/internal/bot"
	"github.com/Wuchieh/discord-bot-template/internal/bot/handler/reaction_role"
	"github.com/Wuchieh/discord-bot-template/internal/config"
	"github.com/Wuchieh/discord-bot-template/internal/database"
	"github.com/Wuchieh/discord-bot-template/internal/model"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

const (
	configFile = "config.yaml"
)

func initConfig() error {
	file, err := os.ReadFile(configFile)
	if err != nil {
		// 找不到設定檔, 創建設定檔
		if errors.Is(err, os.ErrNotExist) {
			cfg := config.GetDefault()
			config.Set(cfg)

			marshal, err := yaml.Marshal(cfg)
			if err != nil {
				return errors.New("配置檔案創建失敗: " + err.Error())
			}

			// 將 YAML 以段落分割，加上空行
			blocks := strings.Split(string(marshal), "\n")
			var result []string
			for i, line := range blocks {
				result = append(result, line)
				// 若下一行是頂層欄位（非縮排開頭）且不是最後一行，就插入空行
				if i+1 < len(blocks) && !strings.HasPrefix(blocks[i+1], " ") && blocks[i+1] != "" {
					result = append(result, "")
				}
			}

			err = os.WriteFile(configFile, []byte(strings.Join(result, "\n")), 0644)
			if err != nil {
				return errors.New("配置檔案創建失敗: " + err.Error())
			}

			fmt.Println("配置檔案創建成功")

			return nil
		}
		return err
	}

	var cfg config.Config

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	config.Set(cfg)

	return nil
}

func initDatabase() error {
	if err := database.Setup(config.Get().DB); err != nil {
		return err
	}

	return database.GetDB().AutoMigrate(model.GetAll()...)
}

func Start() {
	if err := initConfig(); err != nil {
		fmt.Println("配置初始化失敗:", err)
		return
	}

	if err := initDatabase(); err != nil {
		fmt.Println("資料庫初始化失敗:", err)
		return
	}

	reaction_role.Setup(config.Get().ReactionRole)

	defer func() {
		if err := database.Close(); err != nil {
			fmt.Println("資料庫關閉失敗:", err)
		}
	}()

	bot.Start(config.Get().Token)
}

package bootstrap

import (
	"errors"
	"fmt"
	"github.com/Wuchieh/discord-bot-template/internal/bot"
	"github.com/Wuchieh/discord-bot-template/internal/config"
	"github.com/Wuchieh/discord-bot-template/internal/database"
	"github.com/Wuchieh/discord-bot-template/internal/model"
	"gopkg.in/yaml.v3"
	"os"
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

			err = os.WriteFile(configFile, marshal, 0644)
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
	if err := database.Init(); err != nil {
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

	defer func() {
		if err := database.Close(); err != nil {
			fmt.Println("資料庫關閉失敗:", err)
		}
	}()

	bot.Start()
}

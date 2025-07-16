package bot

import (
	"fmt"
	"github.com/Wuchieh/candy-house-bot/internal/bot/handler"
	"github.com/Wuchieh/candy-house-bot/internal/config"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func Start() {
	token := config.Get().Token
	if token == "" {
		fmt.Println("請配置機器人 Token")
		return
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Discord session 創建失敗:", err)
		return
	}

	handler.SetHandler(dg)

	dg.Identify.Intents = discordgo.IntentsAll

	err = dg.Open()
	if err != nil {
		fmt.Println("連線失敗:", err)
		return
	}

	fmt.Printf("%s 已啟動,  按下 CTRL-C 關閉機器人\n", dg.State.User.String())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt, os.Kill)
	<-sc

	if err := dg.Close(); err != nil {
		fmt.Println("關閉連線失敗:", err)
	}
}

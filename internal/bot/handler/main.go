package handler

import "github.com/bwmarrin/discordgo"

var (
	handlers      []any
	onOpened      []func(s *discordgo.Session)
	onBeforeClose []func(s *discordgo.Session)
)

// AddHandler 添加事件 請執行於 init 中
func AddHandler(handler ...any) {
	handlers = append(handlers, handler...)
}

// OnOpened 添加 hook  請執行於 init 中
func OnOpened(fn func(s *discordgo.Session)) {
	onOpened = append(onOpened, fn)
}

// OnBeforeClose 添加 hook  請執行於 init 中
func OnBeforeClose(fn func(s *discordgo.Session)) {
	onBeforeClose = append(onBeforeClose, fn)
}

func SetHandler(s *discordgo.Session) {
	for _, handler := range handlers {
		s.AddHandler(handler)
	}
}

func RunOnOpened(s *discordgo.Session) {
	for _, fn := range onOpened {
		fn(s)
	}
}

func RunOnBeforeClose(s *discordgo.Session) {
	for _, fn := range onBeforeClose {
		fn(s)
	}
}

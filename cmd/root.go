package cmd

import (
	"github.com/Wuchieh/discord-bot-template/internal/bootstrap"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	baseName = filepath.Base(os.Args[0])
)

var rootCmd = &cobra.Command{
	Use: baseName,
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Start()
	},
}

func Execute() {
	_ = rootCmd.Execute()
}

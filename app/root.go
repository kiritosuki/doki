package app

import (
	"github.com/kiritosuki/doki/internal/container"
	"github.com/spf13/cobra"
)

// RootCmd 是所有命令的根命令
var RootCmd = &cobra.Command{
	Use:   "doki",
	Short: "doki is the root of all command",
}

func init() {
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(container.InitCmd)
}

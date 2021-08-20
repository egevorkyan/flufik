package main

import (
	"github.com/egevorkyan/flufik/internal/command"
	"github.com/egevorkyan/flufik/pkg/logging"
)

func main() {
	rootCmd := command.InitialiseRootCmd()
	if err := rootCmd.Command.Execute(); err != nil {
		logging.ErrorHandler("main command execution failed: ", err)
	}
}

package main

import (
	"github.com/egevorkyan/flufik/internal/command"
)

func main() {
	rootCmd := command.InitialiseRootCmd()
	if err := rootCmd.Command.Execute(); err != nil {
		panic(err)
	}
}

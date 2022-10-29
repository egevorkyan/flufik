package command

import (
	"github.com/spf13/cobra"
)

type PushFlufikCommand struct {
	command *cobra.Command
}

func NewFlufikPushCommand() *PushFlufikCommand {
	c := &PushFlufikCommand{
		command: &cobra.Command{
			Use:   "push",
			Short: "any rpm or deb packages to repositories like nexus3, jfrog or generic",
		},
	}
	jrepo := NewFlufikPushJfrogCommand()
	nxrepo := NewFlufikPushNexusCommand()
	f := NewFlufikPushRepoFlufikCommand()
	c.command.AddCommand(
		jrepo.command,
		nxrepo.command,
		f.command,
	)
	return c
}

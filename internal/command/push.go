package command

import (
	"github.com/spf13/cobra"
)

type PushFlufikCommand struct {
	command      *cobra.Command
	provider     string //short - w
	repoUser     string //short - u
	repoPwd      string //short - p
	repoUrl      string //short - l
	packageName  string //short - b
	path         string //short - m
	distribution string //short - d
	component    string //short - c
	architecture string //short - a
	nxcomponent  string //short - n
	repository   string //short - r
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

package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/jfrog"
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
}

func NewFlufikPushCommand() *PushFlufikCommand {
	c := &PushFlufikCommand{
		command: &cobra.Command{
			Use:   "push",
			Short: "pushes any rpm to repository",
		},
	}
	c.command.Flags().StringVarP(&c.provider, "provider", "w", "", "jfrog|nexus|generic")
	c.command.Flags().StringVarP(&c.repoUser, "user", "u", ".", "repository user (must have permission to upload packages)")
	c.command.Flags().StringVarP(&c.repoPwd, "password", "p", "", "repository password")
	c.command.Flags().StringVarP(&c.repoUrl, "url", "l", "", "repository url")
	c.command.Flags().StringVarP(&c.packageName, "package", "b", "", "package name for push")
	c.command.Flags().StringVarP(&c.path, "path", "m", "", "path from where take package")
	c.command.Flags().StringVarP(&c.distribution, "dist", "d", "", "only required for deb packages to push")
	c.command.Flags().StringVarP(&c.component, "component", "c", "main", "only requires for deb packages to push")
	c.command.Flags().StringVarP(&c.architecture, "arch", "a", "", "architecture example: for deb amd64, for rpm x86_64")
	c.command.Run = c.Run
	return c
}

func (c *PushFlufikCommand) Run(command *cobra.Command, args []string) {
	if c.provider == "jfrog" {
		push := jfrog.NewUpload(c.repoUser, c.repoPwd, c.repoUrl, c.packageName, c.path, c.distribution, c.component, c.architecture)
		if err := push.FlufikJFrogUpload(); err != nil {
			logging.ErrorHandler("failure occured during package upload: ", err)
		}
	} else {
		logging.ErrorHandler("provider not provided", fmt.Errorf("provider is empty: %s", c.provider))
	}
}

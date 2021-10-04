package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/jfrog"
	"github.com/egevorkyan/flufik/pkg/plugins/nexus"
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
			Short: "any rpm or deb packages to repositories like nexus3 and jfrog",
		},
	}
	c.command.Flags().StringVarP(&c.provider, "provider", "w", "", "jfrog|nexus|generic")
	c.command.Flags().StringVarP(&c.repoUser, "user", "u", ".", "repository user (must have permission to upload packages)")
	c.command.Flags().StringVarP(&c.repoPwd, "password", "p", "", "repository password")
	c.command.Flags().StringVarP(&c.repoUrl, "url", "l", "", "repository url")
	c.command.Flags().StringVarP(&c.packageName, "package", "b", "", "package name for push")
	c.command.Flags().StringVarP(&c.path, "path", "m", core.FlufikOutputHome(), "path from where take package")
	c.command.Flags().StringVarP(&c.distribution, "dist", "d", "", "only required for deb packages to push")
	c.command.Flags().StringVarP(&c.component, "component", "c", "main", "only requires for deb packages to push")
	c.command.Flags().StringVarP(&c.nxcomponent, "nxcomponent", "n", "", "Nexus components - apt or yum")
	c.command.Flags().StringVarP(&c.repository, "repository", "r", "", "repository name for apt or yum")
	c.command.Flags().StringVarP(&c.architecture, "arch", "a", "", "architecture example: for deb amd64, for rpm x86_64")
	c.command.Run = c.Run
	return c
}

func (c *PushFlufikCommand) Run(command *cobra.Command, args []string) {
	if c.provider == "jfrog" {
		if c.repoUser == "" || c.repoPwd == "" || c.repoUrl == "" || c.packageName == "" || c.distribution == "" || c.component == "" || c.architecture == "" || c.repository == "" {
			logging.ErrorHandler("Warning: ", fmt.Errorf("Required arguments are missing, pushing to jfrog interrupted"))
		} else {
			push := jfrog.NewUpload(c.repoUser, c.repoPwd, c.repoUrl, c.packageName, c.path, c.distribution, c.component, c.architecture, c.repository)
			if err := push.FlufikJFrogUpload(); err != nil {
				logging.ErrorHandler("failure occured during package upload: ", err)
			}
		}
	} else if c.provider == "nexus" {
		if c.repoUser == "" || c.repoPwd == "" || c.repoUrl == "" || c.packageName == "" || c.nxcomponent == "" || c.repository == "" {
			logging.ErrorHandler("Warning: ", fmt.Errorf("Required arguments are missing, pushing to nexus interrupted"))
		} else {
			fnx := nexus.NewNexusUpload(c.repoUser, c.repoPwd, c.repoUrl, c.packageName, c.path, c.nxcomponent, c.repository)
			if err := fnx.FlufikNexusUpload(); err != nil {
				logging.ErrorHandler("Failure: ", err)
			}
		}
	} else {
		logging.ErrorHandler("provider not provided", fmt.Errorf("provider is empty: %s", c.provider))
	}
}

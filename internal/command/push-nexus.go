package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/egevorkyan/flufik/pkg/plugins/nexus"
	"github.com/spf13/cobra"
	"os"
)

type PushNexusFlufikCommand struct {
	command     *cobra.Command
	repoUser    string //short - u ENV NEXUS_REPO_USER
	repoPwd     string //short - p ENV NEXUS_REPO_PWD
	repoUrl     string //short - l ENV NEXUS_REPO_URL
	packageName string //short - b
	packagePath string //short - m
	component   string //short - c
	repository  string //short - r
}

func NewFlufikPushNexusCommand() *PushNexusFlufikCommand {
	c := &PushNexusFlufikCommand{
		command: &cobra.Command{
			Use:   "nexus",
			Short: "push package to nexus repository",
		},
	}
	c.command.Flags().StringVarP(&c.repoUser, "username", "u", os.Getenv("NEXUS_REPO_USER"), "nexus user, possible to get from environment variable NEXUS_REPO_USER")
	c.command.Flags().StringVarP(&c.repoPwd, "password", "p", os.Getenv("NEXUS_REPO_PWD"), "nexus password, possible to get from environment variable NEXUS_REPO_PWD")
	c.command.Flags().StringVarP(&c.repoUrl, "url", "l", os.Getenv("NEXUS_REPO_URL"), "nexus url, possible to get from environment variable NEXUS_REPO_URL")
	c.command.Flags().StringVarP(&c.packageName, "package", "b", "", "package name | example.deb or example.rpm")
	c.command.Flags().StringVarP(&c.packagePath, "package-path", "m", core.FlufikOutputHome(), "location where example.deb or example.rpm located, default is <current user>/.flufik/output/")
	c.command.Flags().StringVarP(&c.component, "component", "c", "main", "component, example main test dev, default is main")
	c.command.Flags().StringVarP(&c.repository, "repo-name", "r", "", "nexus repository name")
	c.command.Run = c.Run
	return c
}

func (c *PushNexusFlufikCommand) Run(command *cobra.Command, args []string) {
	if c.repoUser == "" || c.repoPwd == "" || c.repoUrl == "" || c.packageName == "" || c.component == "" || c.repository == "" {
		logger.RaiseErr("Warning: required arguments are missing, pushing to nexus interrupted")
	} else {
		fnx := nexus.NewNexusUpload(c.repoUser, c.repoPwd, c.repoUrl, c.packageName, c.packagePath, c.component, c.repository)
		if err := fnx.FlufikNexusUpload(); err != nil {
			logger.RaiseErr("failed during push to nexus repository", err)
		}
		logger.InfoLog("successfully pushed to nexus repository")
	}
}

package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/egevorkyan/flufik/pkg/plugins/jfrog"
	"github.com/spf13/cobra"
	"os"
)

type PushJfrogFlufikCommand struct {
	command      *cobra.Command
	repoUser     string //short - u ENV JFROG_REPO_USER
	repoPwd      string //short - p ENV JFROG_REPO_PWD
	repoUrl      string //short - l ENV JFROG_REPO_URL
	packageName  string //short - b
	packagePath  string //short - m
	distribution string //short - d
	component    string //short - c
	repository   string //short - r
}

func NewFlufikPushJfrogCommand() *PushJfrogFlufikCommand {
	c := &PushJfrogFlufikCommand{
		command: &cobra.Command{
			Use:   "jfrog",
			Short: "push package to jfrog repository",
		},
	}
	c.command.Flags().StringVarP(&c.repoUser, "username", "u", os.Getenv("JFROG_REPO_USER"), "jfrog user, possible to get from environment variable JFROG_REPO_USER")
	c.command.Flags().StringVarP(&c.repoPwd, "password", "p", os.Getenv("JFROG_REPO_PWD"), "jfrog password, possible to get from environment variable JFROG_REPO_PWD")
	c.command.Flags().StringVarP(&c.repoUrl, "url", "l", os.Getenv("JFROG_REPO_URL"), "jfrog url, possible to get from environment variable JFROG_REPO_URL")
	c.command.Flags().StringVarP(&c.packageName, "package", "b", "", "package name | example.deb or example.rpm")
	c.command.Flags().StringVarP(&c.packagePath, "package-path", "m", core.FlufikOutputHome(), "location where example.deb or example.rpm located, default is <current user>/.flufik/output/")
	c.command.Flags().StringVarP(&c.distribution, "distribution", "d", "", "distribution, example focal, flufik")
	c.command.Flags().StringVarP(&c.component, "component", "c", "main", "component, example main test dev, default is main")
	c.command.Flags().StringVarP(&c.repository, "repo-name", "r", "", "jfrog repository name")
	c.command.Run = c.Run
	return c
}

func (c *PushJfrogFlufikCommand) Run(command *cobra.Command, args []string) {
	arch := core.CheckArch(c.packageName)
	if len(arch) > 0 {
		if c.repoUser == "" || c.repoPwd == "" || c.repoUrl == "" || c.packageName == "" || c.distribution == "" || c.component == "" || c.repository == "" {
			logger.RaiseErr("required arguments are missing, pushing to jfrog interrupted")
		} else {
			push := jfrog.NewUpload(c.repoUser, c.repoPwd, c.repoUrl, c.packageName, c.packagePath, c.distribution, c.component, arch, c.repository)
			if err := push.FlufikJFrogUpload(); err != nil {
				logger.RaiseErr("failure occurred during package upload", err)
			}
			logger.InfoLog("successfully pushed to jfrog repository")
		}
	} else {
		logger.RaiseWarn("Warning: package name is not based on official naming convention")
	}
}

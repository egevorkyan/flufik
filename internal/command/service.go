package command

import "github.com/spf13/cobra"

type ServiceFlufikCommand struct {
	command       *cobra.Command
	serviceConfig string
}

func NewFlufikServiceCommand() *ServiceFlufikCommand {
	c := &ServiceFlufikCommand{
		command: &cobra.Command{
			Use:   "service",
			Short: "starts service",
		},
	}
	c.command.Run = c.Run
	return c
}

func (c *ServiceFlufikCommand) Run(command *cobra.Command, args []string) {

}

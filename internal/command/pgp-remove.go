package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
)

type PgpFlufikRemoveCommand struct {
	command    *cobra.Command
	pgpKeyName string
}

func NewPgpFlufikRemoveCommand() *PgpFlufikRemoveCommand {
	c := &PgpFlufikRemoveCommand{
		command: &cobra.Command{
			Use:   "remove",
			Short: "remove pgp key",
		},
	}
	c.command.Flags().StringVarP(&c.pgpKeyName, "name", "n", "", "Provide key name to save")
	c.command.Run = c.Run
	return c
}

func (c *PgpFlufikRemoveCommand) Run(command *cobra.Command, args []string) {
	if err := crypto.RemovePgpKeyFromDB(c.pgpKeyName); err != nil {
		logging.ErrorHandler("info: ", err)
	} else {
		logging.ErrorHandler("info: ", fmt.Errorf("successfully removed"))
	}
}

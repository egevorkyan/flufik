package command

import (
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logger"
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
	p := pgp.NewImportPGP()
	if err := p.RemovePgpKeyFromDB(c.pgpKeyName); err != nil {
		logger.RaiseErr("pgp key removal process failed", err)
	}
	logger.InfoLog("successfully removed")
}

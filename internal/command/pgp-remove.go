package command

import (
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
	"os"
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
	logger := logging.GetLogger()
	debuging := os.Getenv("FLUFIK_DEBUG")
	if debuging == "1" {
		logger.Info("remove pgp key")
	}
	p := pgp.NewImportPGP(logger, debuging)
	if err := p.RemovePgpKeyFromDB(c.pgpKeyName); err != nil {
		logger.Errorf("pgp key removal process failed: %v", err)
	} else {
		logger.Info("pgp key successfully removed")
	}
}

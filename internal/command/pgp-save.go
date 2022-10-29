package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/spf13/cobra"
)

type PgpFlufikDumpCommand struct {
	command    *cobra.Command
	pgpKeyName string
	location   string
}

func NewPgpFlufikSaveCommand() *PgpFlufikDumpCommand {
	c := &PgpFlufikDumpCommand{
		command: &cobra.Command{
			Use:   "export",
			Short: "export pgp key to file if required, passphrase will be printed on screen",
		},
	}
	c.command.Flags().StringVarP(&c.pgpKeyName, "name", "n", "", "Provide key name to save")
	c.command.Flags().StringVarP(&c.location, "path", "p", core.FlufikCurrentDir(), "path where to save keys")
	c.command.Run = c.Run
	return c
}

func (c *PgpFlufikDumpCommand) Run(command *cobra.Command, args []string) {
	p := pgp.NewImportPGP()
	if err := p.SavePgpKeyToFile(c.pgpKeyName, c.location); err != nil {
		logger.RaiseErr("error occurred during export pgp key", err)
	}
	logger.InfoLog("successfully saved")
}

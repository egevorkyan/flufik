package command

import (
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/spf13/cobra"
)

type PgpFlufikGenerateCommand struct {
	command *cobra.Command
	name    string
	email   string
	comment string
	keyType string
	bits    int
}

func NewFlufikPgpGenerateCommand() *PgpFlufikGenerateCommand {
	c := &PgpFlufikGenerateCommand{
		command: &cobra.Command{
			Use:   "gen",
			Short: "generates pgp key with passphrase",
		},
	}
	c.command.Flags().StringVarP(&c.name, "name", "n", "", "pgp key name")
	c.command.Flags().StringVarP(&c.email, "email", "e", ".", "email address")
	c.command.Flags().StringVarP(&c.comment, "comment", "c", "", "pgp comment")
	c.command.Flags().StringVarP(&c.keyType, "key-type", "k", "", "default key type is rsa. possible types: rsa|x25519. In case of x25519 bits values is not required")
	c.command.Flags().IntVarP(&c.bits, "bits", "b", 0, "pgp key bits")
	c.command.Run = c.Run
	return c
}

func (c *PgpFlufikGenerateCommand) Run(command *cobra.Command, args []string) {
	p := pgp.NewPGP(c.name, c.email, c.comment, c.keyType, c.bits)
	if err := p.GeneratePgpKey(); err != nil {
		logger.RaiseErr("", err)
	}
	logger.InfoLog("successfully generated")
}

package command

import (
	"fmt"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
)

type FlufikPgpImportCommand struct {
	command    *cobra.Command
	privateKey string
	publicKey  string
	passPhrase string
	name       string
}

func NewFlufikImportPgpKey() *FlufikPgpImportCommand {
	c := &FlufikPgpImportCommand{
		command: &cobra.Command{
			Use:   "import",
			Short: "importing pgp keys with passphrase",
		},
	}
	c.command.Flags().StringVarP(&c.name, "name", "n", "", "Key name")
	c.command.Flags().StringVarP(&c.privateKey, "private", "p", "", "Private Key Path")
	c.command.Flags().StringVarP(&c.publicKey, "public", "c", "", "Public Key Path")
	c.command.Flags().StringVarP(&c.passPhrase, "passphrase", "s", "", "Pricate Key Passphrase")
	c.command.Run = c.Run
	return c
}

func (c *FlufikPgpImportCommand) Run(command *cobra.Command, args []string) {
	if c.passPhrase == "" {
		logging.ErrorHandler("message: ", fmt.Errorf("only pgp key with passphrase is accepted"))
	}
	if err := crypto.ImportPgpKeys(c.name, c.privateKey, c.publicKey, c.passPhrase); err != nil {
		logging.ErrorHandler("fatal: %w", err)
	}
}

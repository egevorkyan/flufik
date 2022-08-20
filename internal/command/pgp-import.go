package command

import (
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
	"os"
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
	logger := logging.GetLogger()
	debuging := os.Getenv("FLUFIK_DEBUG")
	if debuging == "1" {
		logger.Info("import pgp key")
	}
	if c.passPhrase == "" {
		logger.Errorf("only pgp key with passphrase is accepted")
	}
	p := pgp.NewImportPGP(logger, debuging)
	if err := p.ImportPgpKeys(c.name, c.privateKey, c.publicKey, c.passPhrase); err != nil {
		logger.Errorf("fatal: %v", err)
	}
}

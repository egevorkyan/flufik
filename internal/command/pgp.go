package command

import "github.com/spf13/cobra"

type FlufikPgpCmd struct {
	commnad *cobra.Command
}

func NewFlufikPgp() *FlufikPgpCmd {
	c := &FlufikPgpCmd{
		commnad: &cobra.Command{
			Use:   "pgp",
			Short: "pgp releated menu",
		},
	}
	pgpGen := NewFlufikPgpGenerateCommand()
	pgpSave := NewPgpFlufikSaveCommand()
	pgpImport := NewFlufikImportPgpKey()
	c.commnad.AddCommand(
		pgpGen.command,
		pgpSave.command,
		pgpImport.command,
	)
	return c
}

package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/spf13/cobra"
)

type RootFlufikCommand struct {
	Command *cobra.Command
}

func NewFlufikRootCommand() *RootFlufikCommand {
	c := &RootFlufikCommand{
		Command: &cobra.Command{
			Use:   "flufik",
			Short: "Flufik: CLI tool for building awesome rpm and deb packages",
			Long: `
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
|                                                        |
|                /,,,,\_____________/,,,,\               |
|               |,(  )/,,,,,,,,,,,,,\(  ),|              |
|                \__,,,,___,,,,,___,,,,__/               |
|                  /,,,/(')\,,,/(')\,,,\                 |
|                 |,,,,___ _____ ___,,,,|                |
|                 |,,,/   \\o_o//   \,,,|                |
|                 |,,|       |       |,,|                |
|                 |,,|   \__/|\__/   |,,|                |
|                  \,,\     \_/     /,,/                 |
|                   \__\___________/__/                  |
|     ________________/,,,,,,,,,,,,,\________________    |
|    / \,,,,,,,,,,,,,,,,___________,,,,,,,,,,,,,,,,/ \   |
|   (   ),,,,,,,,,,,,,,/           \,,,,,,,,,,,,,,(   )  |
|    \_/____________,,/             \,,____________\_/   |
|                  /,/               \,\                 |
|                 |,|   I am Flufik   |,|                |
|                 |,|  ready to pack  |,|                |
|                 |,|  apps for Linux |,|                |
|                 |,|                 |,|                |
|                  \,\       O       /,/                 |
|                  /,,\_____________/,,\                 |
|                 /,,,,,,,,,,,,,,,,,,,,,\                |
|                /,,,,,,,,_______,,,,,,,,\               |
|               /,,,,,,,,/       \,,,,,,,,\              |
|              /,,,,,,, /         \,,,,,,,,\             |
|             /_____,,,/           \,,,_____\            |
|            //     \,/             \,/     \\           |
|            \\_____//               \\_____//           |
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
`,
			Version: core.VERSION,
		},
	}
	c.Command.SetVersionTemplate("flufik packager version {{.Version}}\n")
	cobra.OnInitialize(c.initConfig)
	return c
}

func (c *RootFlufikCommand) initConfig() {

}

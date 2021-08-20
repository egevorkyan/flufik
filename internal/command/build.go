package command

import (
	"bufio"
	"github.com/egevorkyan/flufik/internal/flufikbuilder"
	"github.com/egevorkyan/flufik/internal/flufikinfo"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/spf13/cobra"
	"os"
	"path"
)

type BuildFlufikCommand struct {
	command         *cobra.Command
	buildPack       string
	srcDir          string
	buildConfigPath string
}

func NewFlufikBuildCommand() *BuildFlufikCommand {
	c := &BuildFlufikCommand{
		command: &cobra.Command{
			Use:   "build",
			Short: "builds deployment rpm or deb or both packages",
		},
	}
	c.command.Flags().StringVarP(&c.buildPack, "package", "p", "", "used to identify what type of package to build: values rpm|deb")
	c.command.Flags().StringVarP(&c.srcDir, "source-directory", "s", ".", "source directory default is current location")
	c.command.Flags().StringVarP(&c.buildConfigPath, "configuration-file", "c", "config.yaml", "configuration file used during build, default is current location config.yaml")
	c.command.Run = c.Run
	return c
}

func (c *BuildFlufikCommand) Run(command *cobra.Command, args []string) {
	pkgInfoLoader, err := flufikinfo.LoadPackageInfo(c.buildConfigPath, c.srcDir)
	if err != nil {
		logging.ErrorHandler("can't load configuration file error: ", err)
	}
	switch c.buildPack {
	case "rpm":
		if err = buildFlufikPackage(flufikbuilder.NewFlufikRpmBuilder(pkgInfoLoader), c.srcDir); err != nil {
			logging.ErrorHandler("rpm package not build error: ", err)
		}
	case "deb":
		if err = buildFlufikPackage(flufikbuilder.NewFlufikDebBuilder(pkgInfoLoader), c.srcDir); err != nil {
			logging.ErrorHandler("deb package not build error: ", err)
		}
	}
}

func buildFlufikPackage(flufikBuilder flufikbuilder.FlufikPackageBuilder, directory string) error {
	var pkgFile *os.File
	if pkgPath, err := flufikBuilder.FileName(); err == nil {
		pkgPath = path.Join(directory, pkgPath)

		if pkgFile, err = os.Create(pkgPath); err != nil {
			return err
		}
	} else {
		return err
	}
	defer func() { _ = pkgFile.Close() }()

	pkgWriter := bufio.NewWriter(pkgFile)

	if err := flufikBuilder.Build(pkgWriter); err != nil {
		return err
	}

	if err := pkgWriter.Flush(); err != nil {
		return err
	}
	return nil
}

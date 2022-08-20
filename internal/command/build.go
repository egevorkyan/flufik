package command

import (
	"bufio"
	"github.com/egevorkyan/flufik/core"
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
	destDir         string
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
	c.command.Flags().StringVarP(&c.destDir, "destination-directory", "d", core.FlufikOutputHome(), "output directory default is current user ~/.flufik/output")
	c.command.Flags().StringVarP(&c.buildConfigPath, "configuration-file", "c", "config.yaml", "configuration file used during build, default is current location config.yaml")
	c.command.Run = c.Run
	return c
}

func (c *BuildFlufikCommand) Run(command *cobra.Command, args []string) {
	logger := logging.GetLogger()
	debuging := os.Getenv("FLUFIK_DEBUG")
	if debuging == "1" {
		logger.Info("build packages")
	}
	pkgInfoLoader, err := flufikinfo.LoadPackageInfo(c.buildConfigPath, c.destDir)
	if err != nil {
		logger.Errorf("can't load configuration file error: %v", err)
	}
	switch c.buildPack {
	case "rpm":
		if err = buildFlufikPackage(flufikbuilder.NewFlufikRpmBuilder(pkgInfoLoader, logger, debuging), c.destDir); err != nil {
			logger.Errorf("rpm package not build error: %v", err)
		}
	case "deb":
		if err = buildFlufikPackage(flufikbuilder.NewFlufikDebBuilder(pkgInfoLoader, logger, debuging), c.destDir); err != nil {
			logger.Errorf("deb package not build error: %v", err)
		}
	}
}

func buildFlufikPackage(flufikBuilder flufikbuilder.FlufikPackageBuilder, directory string) error {
	var pkgFile *os.File
	var dst string
	pkgPath, err := flufikBuilder.FileName()
	if err != nil {
		return err
	}
	if path.IsAbs(directory) {
		dst = directory
	} else {
		dst = path.Join(core.FlufikCurrentDir(), directory)
	}
	if _, err = os.Stat(dst); os.IsNotExist(err) {
		err = os.MkdirAll(dst, 0755)
		if err != nil {
			return err
		}
	}
	p := core.FlufikPkgFilePath(pkgPath, dst)
	if err != nil {
		return err
	}
	if pkgFile, err = os.Create(p); err != nil {
		return err
	}
	defer func() { _ = pkgFile.Close() }()

	pkgWriter := bufio.NewWriter(pkgFile)

	if err = flufikBuilder.Build(pkgWriter); err != nil {
		return err
	}

	if err = pkgWriter.Flush(); err != nil {
		return err
	}
	return nil
}

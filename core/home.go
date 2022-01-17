package core

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	FLUFIKHOME             = ".flufik"
	FLUFIKKEYSDIR          = "keys"
	FLUFIKKEYSDB           = "keys"
	FLUFIKLOGGINGDIR       = "logs"
	FLUFIKPKGCONFIGDIR     = "configs" //yaml configuration file or template, to build based on that package
	FLUFIKPKGOUTPUTDIR     = "output"
	FLUFIKLOGGINGFILE      = "all.log"
	FLUFIKSERVICECONFIGDIR = "service"
)

func Home() string {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return u.HomeDir
}

func FlufikHome() string {
	return filepath.Join(Home(), FLUFIKHOME)
}

func FlufikKeysHome() string {
	return filepath.Join(FlufikHome(), FLUFIKKEYSDIR)
}

func FlufikOutputHome() string {
	return filepath.Join(FlufikHome(), FLUFIKPKGOUTPUTDIR)
}

func FlufikLoggingHome() string {
	return filepath.Join(FlufikHome(), FLUFIKLOGGINGDIR)
}

func FlufikConfigurationHome() string {
	return filepath.Join(FlufikHome(), FLUFIKPKGCONFIGDIR)
}

func FlufikLoggingFilePath() string {
	return filepath.Join(FlufikLoggingHome(), FLUFIKLOGGINGFILE)
}

func FlufikKeyFileName(private, public, extention string) (string, string) {
	return filepath.Join(FlufikKeysHome(), fmt.Sprintf("%s.%s", private, extention)), filepath.Join(FlufikKeysHome(), fmt.Sprintf("%s.%s", public, extention))
}

func FlufikKeyFilePath(name string) string {
	return filepath.Join(FlufikKeysHome(), name)
}

func FlufikKeyDbPath() string {
	return filepath.Join(FlufikKeysHome(), FLUFIKKEYSDB)
}

func FlufikPkgFilePath(pkg, path string) string {
	return filepath.Join(path, pkg)
}

func FlufikServiceConfigurationHome() string {
	return filepath.Join(FlufikHome(), FLUFIKSERVICECONFIGDIR)
}

func FlufikCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

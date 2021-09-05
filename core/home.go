package core

import (
	"log"
	"os/user"
	"path/filepath"
)

const FLUFIKHOME = ".flufik"
const FLUFIKKEYSDIR = "keys"
const FLUFIKLOGGINGDIR = "logs"
const FLUFIKPKGCONFIGDIR = "configs" //yaml configuration file or template, to build based on that package

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

func FlufikLoggingHome() string {
	return filepath.Join(FlufikHome(), FLUFIKLOGGINGDIR)
}

func FlufikConfigurationHome() string {
	return filepath.Join(FlufikHome(), FLUFIKPKGCONFIGDIR)
}

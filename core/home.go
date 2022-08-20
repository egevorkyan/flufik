package core

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	FLUFIKHOME             = ".flufik"
	FLUFIKKEYSDIR          = "keys"
	USERDB                 = "user.db"
	PGPDB                  = "pgp.db"
	FLUFIKLOGGINGDIR       = "logs"
	FLUFIKPKGCONFIGDIR     = "configs" //yaml configuration file or template, to build based on that package
	FLUFIKPKGOUTPUTDIR     = "output"
	FLUFIKLOGGINGFILE      = "all.log"
	FLUFIKSERVICECONFIGDIR = "service"
	FLUFIKSERVICEWEB       = "flufikweb"
	FLUFIKROOTPATH         = "/opt/flufik"
	FLUFIKNOSQLDB          = "flufikdb"
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

func FlufikUserDbPath() string {
	return filepath.Join(FlufikKeysHome(), USERDB)
}

func FlufikPgpDbPath() string {
	return filepath.Join(FlufikKeysHome(), PGPDB)
}

func FlufikNoSqlDbPath() string {
	return filepath.Join(FlufikKeysHome(), FLUFIKNOSQLDB)
}

func FlufikPkgFilePath(pkg, dirPath string) string {
	return filepath.Join(dirPath, pkg)
}

func FlufikServiceConfigurationHome() string {
	return filepath.Join(FlufikHome(), FLUFIKSERVICECONFIGDIR)
}

func FlufikCurrentDir() string {
	currentDir, _ := os.Getwd()
	return currentDir
}

func FlufikRootHome() string {
	return FLUFIKROOTPATH
}

func FlufikServiceWebHome(rootPath string) string {
	return filepath.Join(rootPath, FLUFIKSERVICEWEB)
}

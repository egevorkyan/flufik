package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/users"
	"log"
	"os"
)

func init() {
	flufikHomeInit()
}

func flufikHomeInit() {
	logger := logging.GetLogger()
	debuging := os.Getenv("FLUFIK_DEBUG")
	if debuging == "1" {
		logger.Info("initialize flufik")
	}
	_, err := os.Stat(core.FlufikHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikHome(), os.ModePerm); err != nil {
			logger.Fatalf("can not create initial directories: %v", err)
		}
	}

	_, err = os.Stat(core.FlufikLoggingHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikLoggingHome(), os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	_, err = os.Stat(core.FlufikKeysHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikKeysHome(), os.ModePerm); err != nil {
			logger.Errorf("can not create keys folder: %v", err)
		}
	}

	_, err = os.Stat(core.FlufikConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikConfigurationHome(), os.ModePerm); err != nil {
			logger.Errorf("can not create configuration folder: %v", err)
		}
	}

	_, err = os.Stat(core.FlufikOutputHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikOutputHome(), os.ModePerm); err != nil {
			logger.Errorf("can not create output folder: %v", err)
		}
	}

	_, err = os.Stat(core.FlufikServiceConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikServiceConfigurationHome(), os.ModePerm); err != nil {
			logger.Errorf("can not create output folder: %v", err)
		}
	}

	if err = initializePgpKey(logger, debuging); err != nil {
		logger.Errorf("initialize default pgp key: %v", err)
	}

	if err = initializeAdminUser(logger, debuging); err != nil {
		logger.Errorf("initialize admin user: %v", err)
	}

}

func initializeAdminUser(logger *logging.Logger, debugging string) error {
	u := users.NewUser(logger, debugging)
	err := u.CreateUser("admin", "admin")
	if err != nil {
		return err
	}
	err = u.DumpUser("admin", "initial-user.txt")
	if err != nil {
		return err
	}
	return nil
}

func initializePgpKey(logger *logging.Logger, debugging string) error {
	p := pgp.NewPGP("flufik", "", "", "rsa", 4096, logger, debugging)
	err := p.GeneratePgpKey()
	if err != nil {
		return err
	}
	return nil
}

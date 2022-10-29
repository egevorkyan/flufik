package command

import (
	"errors"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto/pgp"
	"github.com/egevorkyan/flufik/pkg/logger"
	"github.com/egevorkyan/flufik/users"
	"os"
)

func init() {
	flufikHomeInit()
}

func flufikHomeInit() {
	_, err := os.Stat(core.FlufikHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikHome(), os.ModePerm); err != nil {
			logger.RaiseErr("can not create initial directories", err)
		}
	}

	_, err = os.Stat(core.FlufikLoggingHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikLoggingHome(), os.ModePerm); err != nil {
			logger.RaiseErr("failed to create initial directories", err)
		}
	}

	_, err = os.Stat(core.FlufikKeysHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikKeysHome(), os.ModePerm); err != nil {
			logger.RaiseErr("can not create keys folder", err)
		}
	}

	_, err = os.Stat(core.FlufikConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikConfigurationHome(), os.ModePerm); err != nil {
			logger.RaiseErr("can not create configuration folder", err)
		}
	}

	_, err = os.Stat(core.FlufikOutputHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikOutputHome(), os.ModePerm); err != nil {
			logger.RaiseErr("can not create output folder", err)
		}
	}

	_, err = os.Stat(core.FlufikServiceConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikServiceConfigurationHome(), os.ModePerm); err != nil {
			logger.RaiseErr("can not create output folder", err)
		}
	}

	_, err = os.Stat(core.FlufikNoSqlDbPath())
	if errors.Is(err, os.ErrNotExist) {

		if err = initializePgpKey(); err != nil {
			logger.RaiseErr("initialize default pgp key", err)
		}

		if err = initializeAdminUser(); err != nil {
			logger.RaiseErr("initialize admin user", err)
		}
	}

}

func initializeAdminUser() error {
	u := users.NewUser()
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

func initializePgpKey() error {
	p := pgp.NewPGP("flufik", "", "", "rsa", 4096)
	err := p.GeneratePgpKey()
	if err != nil {
		return err
	}
	return nil
}

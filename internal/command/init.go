package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
	"github.com/egevorkyan/flufik/users"
	"log"
	"os"
)

func init() {
	flufikHomeInit()
}

func flufikHomeInit() {
	_, err := os.Stat(core.FlufikHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikHome(), os.ModePerm); err != nil {
			log.Fatal(err)
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
			logging.ErrorHandler("can not create keys folder: ", err)
		}
	}

	_, err = os.Stat(core.FlufikConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikConfigurationHome(), os.ModePerm); err != nil {
			logging.ErrorHandler("can not create configuration folder", err)
		}
	}

	_, err = os.Stat(core.FlufikOutputHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikOutputHome(), os.ModePerm); err != nil {
			logging.ErrorHandler("can not create output folder", err)
		}
	}

	_, err = os.Stat(core.FlufikServiceConfigurationHome())
	if os.IsNotExist(err) {
		if err = os.Mkdir(core.FlufikServiceConfigurationHome(), os.ModePerm); err != nil {
			logging.ErrorHandler("can not create output folder", err)
		}
	}

	if err = initDB(); err != nil {
		logging.ErrorHandler("fatal: ", err)
	}

	if err = initializeAdminUser(); err != nil {
		logging.ErrorHandler("warning: ", err)
	}

}

func initDB() error {
	db, err := simpledb.CreateInternalDb(core.FlufikDbPath())
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
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

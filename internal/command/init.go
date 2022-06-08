package command

import (
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/crypto"
	"github.com/egevorkyan/flufik/pkg/logging"
	"github.com/egevorkyan/flufik/pkg/plugins/simpledb"
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

	if err = initializeKeys(); err != nil {
		logging.ErrorHandler("warning: ", err)
	}

}

func initDB() error {
	db := simpledb.NewSimpleDB(core.FlufikDbPath())
	if err := db.CreateTable(core.FLUFIKKEYDBTYPE); err != nil {
		return err
	}
	if err := db.CreateTable(core.FLUFIKAPPDBTYPE); err != nil {
		return err
	}
	db.CloseDb()
	return nil
}

func initializeKeys() error {
	_, err := crypto.GetApiKey()
	if err = crypto.CreateApiKey(); err != nil {
		return err
	}
	return nil
}

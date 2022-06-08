package simpledb

import (
	"database/sql"
	"fmt"
	"github.com/egevorkyan/flufik/pkg/logging"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type SimpleDb struct {
	sdb *sql.DB
}

func (sDB *SimpleDb) CreateTable(typeTable string) error {
	switch typeTable {
	case "kvdb":
		statement, err := sDB.sdb.Prepare(flufikTableSqlSchema)
		if err != nil {
			return err
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}
		// Do not duplicate data in DB
		statement, err = sDB.sdb.Prepare(flufikTableSqlUnique)
		if err != nil {
			log.Fatal(err.Error())
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}
	case "app":
		statement, err := sDB.sdb.Prepare(appTableSqlSchema)
		if err != nil {
			return err
		}
		_, err = statement.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
func (sDB *SimpleDb) Insert(typeTable string, data ...any) error {
	switch typeTable {
	case "kvdb":
		statement, err := sDB.sdb.Prepare(internalFlufikInsertSchema)
		if err != nil {
			return err
		}
		if len(data) == 4 {
			_, err = statement.Exec(data[0], data[1], data[2], data[3])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s", "missing values")
		}
	case "app":
		statement, err := sDB.sdb.Prepare(internalAppInsertSchema)
		if err != nil {
			return err
		}
		if len(data) == 5 {
			_, err = statement.Exec(data[0], data[1], data[2], data[3], data[4])
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s", "missing values")
		}
	}
	return nil
}

func (sDB *SimpleDb) GetKey(searchData string) (Data, error) {
	var d Data
	row, err := sDB.sdb.Query("SELECT key, privateKey, publicKey, token FROM flufikTab WHERE key = ?", searchData)
	if err != nil {
		return Data{}, err
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&d.KeyValue, &d.PrivateKeyValue, &d.PublicKeyValue, &d.TokenValue)
	}
	return d, nil
}

func (sDB *SimpleDb) GetLatestApp(searchData map[string]interface{}) (App, error) {
	var a App
	row, err := sDB.sdb.Query("SELECT appname, appversion, apparch, apposversion, applocation FROM appTab WHERE appname = ? AND apparch = ? AND apposversion = ?",
		searchData["appName"], searchData["appArch"], searchData["appOsVersion"])
	if err != nil {
		return App{}, err
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&a.AppName, &a.AppVersion, &a.AppArch, &a.AppOsVersion, &a.AppLocation)

	}
	return a, nil
}

func (sDB *SimpleDb) GetAppByVersion(searchData map[string]interface{}) (App, error) {
	var a App
	row, err := sDB.sdb.Query("SELECT appname, appversion, apparch, apposversion, applocation FROM appTab WHERE appname = ? AND apparch = ? AND apposversion = ? AND appversion = ?",
		searchData["appName"], searchData["appArch"], searchData["appOsVersion"], searchData["appVersion"])
	if err != nil {
		return App{}, err
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&a.AppName, &a.AppVersion, &a.AppArch, &a.AppOsVersion, &a.AppLocation)

	}
	return a, nil
}

func (sDB *SimpleDb) Delete(key string) error {
	statement, err := sDB.sdb.Prepare(DeleteFlufikSql)
	if err != nil {
		return err
	}
	_, err = statement.Exec(key)
	if err != nil {
		return err
	}
	return nil
}

func (sDB *SimpleDb) DeleteApp(appName string, appVersion string, appArch string, appOsVersion string) error {
	if len(appName) == 0 {
		return fmt.Errorf("%s", "application name must present")
	} else if len(appName) != 0 && len(appVersion) != 0 && len(appArch) != 0 && len(appOsVersion) != 0 {
		statement, err := sDB.sdb.Prepare(DeleteAppSqlByVersionByArchByOs)
		if err != nil {
			return err
		}
		_, err = statement.Exec(appName, appVersion, appArch, appOsVersion)
		if err != nil {
			return err
		}
	} else if len(appName) != 0 && len(appVersion) != 0 && len(appArch) != 0 {
		statement, err := sDB.sdb.Prepare(DeleteAppSqlByVersionByArch)
		if err != nil {
			return err
		}
		_, err = statement.Exec(appName, appVersion, appArch)
		if err != nil {
			return err
		}
	} else if len(appName) != 0 && len(appVersion) != 0 {
		statement, err := sDB.sdb.Prepare(DeleteAppSqlByVersion)
		if err != nil {
			return err
		}
		_, err = statement.Exec(appName, appVersion)
		if err != nil {
			return err
		}
	} else {
		statement, err := sDB.sdb.Prepare(DeleteAppSql)
		if err != nil {
			return err
		}
		_, err = statement.Exec(appName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sDB *SimpleDb) CloseDb() {
	_ = sDB.sdb.Close()
}

func NewSimpleDB(dbPath string) *SimpleDb {
	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			logging.ErrorHandler("fatal: ", err)
		}
		_ = file.Close()
	}
	simpleDb, _ := sql.Open("sqlite", dbPath)
	sdb := &SimpleDb{
		sdb: simpleDb,
	}
	return sdb
}

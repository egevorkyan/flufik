package simpledb

import (
	"database/sql"
	"github.com/egevorkyan/flufik/core"
	"github.com/egevorkyan/flufik/pkg/logging"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type SimpleDb struct {
	sdb *sql.DB
}
type Data struct {
	KeyValue        string
	PrivateKeyValue string
	PublicKeyValue  string
	TokenValue      string
}

func (sDB *SimpleDb) CreateTable() error {
	createFlufikTableSQL := `CREATE TABLE IF NOT EXISTS flufikTab (
		"idKey" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"key" TEXT,
		"privateKey" TEXT,
		"publicKey" TEXT,
		"token" TEXT
	  );` // SQL Statement for Create Table

	statement, err := sDB.sdb.Prepare(createFlufikTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	// Do not duplicate data in DB
	createFlufikTableSQL = `CREATE UNIQUE index IF NOT EXISTS flufik_key on flufikTab(key)`
	statement, err = sDB.sdb.Prepare(createFlufikTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}
func (sDB *SimpleDb) Insert(key string, privateKey string, publicKey string, token string) error {
	insertFlufikSQL := `INSERT OR ignore INTO flufikTab (key, privateKey, publicKey, token) VALUES (?, ?, ?, ?)`
	statement, err := sDB.sdb.Prepare(insertFlufikSQL)
	// Prepare statement. This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(key, privateKey, publicKey, token)
	if err != nil {
		return err
	}
	return nil
}
func (sDB *SimpleDb) Get(key string) (Data, error) {
	var d Data
	row, err := sDB.sdb.Query("SELECT key, privateKey, publicKey, token FROM flufikTab WHERE key = ?", key)
	if err != nil {
		return Data{}, err
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		row.Scan(&d.KeyValue, &d.PrivateKeyValue, &d.PublicKeyValue, &d.TokenValue)

	}
	return d, nil
}

func (sDB *SimpleDb) Delete(key string) error {
	deleteFlufikSQL := `DELETE FROM flufikTab WHERE key = ?`
	statement, err := sDB.sdb.Prepare(deleteFlufikSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(key)
	if err != nil {
		return err
	}
	return nil
}

func (sDB *SimpleDb) CloseDb() {
	_ = sDB.sdb.Close()
}

func NewSimpleDB() *SimpleDb {
	_, err := os.Stat(core.FlufikKeyDbPath())
	if os.IsNotExist(err) {
		file, err := os.Create(core.FlufikKeyDbPath())
		if err != nil {
			logging.ErrorHandler("fatal: ", err)
		}
		_ = file.Close()
	}
	simpleDb, _ := sql.Open("sqlite", core.FlufikKeyDbPath())
	sdb := &SimpleDb{
		sdb: simpleDb,
	}

	return sdb
}

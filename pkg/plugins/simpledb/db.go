package simpledb

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
)

type FluffDb struct {
	db     *sql.DB
	dbpath string
}

type Pgp struct {
	KeyName    string
	PrivateKey string
	PublicKey  string
	PassPhrase string
}
type User struct {
	UserName string
	Password string
	Mode     string
}

//CreateInternalDb - Create internal Database
func CreateInternalDb(path string) (*FluffDb, error) {
	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			return nil, err
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB: %v", err)
	}

	// create database tables
	_, err = db.Exec(sqlCreateTables)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB tables: %v", err)
	}

	// create database indexes
	_, err = db.Exec(sqlCreateIndexes)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB indexes: %v", err)
	}

	return &FluffDb{
		db:     db,
		dbpath: path,
	}, nil
}

// OpenInternalDB opens a internal database SQLite database from file and return a
// pointer to the resulting struct.
func OpenInternalDB(path string) (*FluffDb, error) {
	// open database file
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	return &FluffDb{
		db:     db,
		dbpath: path,
	}, nil
}

func (f *FluffDb) Begin() (*sql.Tx, error) {
	return f.db.Begin()
}

func (f *FluffDb) Close() error {
	if f.db != nil {
		return f.db.Close()
	}
	return nil
}

func (f *FluffDb) InsertPgpKeys(keyName string, privateKey string, publicKey string, passPhrase string) error {
	// insert pgp key
	statement, err := f.db.Prepare(sqlInsertPgpKey)
	if err != nil {
		return err
	}
	defer func(statement *sql.Stmt) {
		err = statement.Close()
		if err != nil {
			return
		}
	}(statement)

	_, err = statement.Exec(keyName, privateKey, publicKey, passPhrase)
	if err != nil {
		return err
	}
	return nil
}

func (f *FluffDb) InsertUsers(userName string, password string, mode string) error {
	// insert  user
	statement, err := f.db.Prepare(sqlInsertUser)
	if err != nil {
		return err
	}
	defer func(statement *sql.Stmt) {
		err = statement.Close()
		if err != nil {
			return
		}
	}(statement)

	_, err = statement.Exec(userName, password, mode)
	if err != nil {
		return err
	}
	return nil
}

func (f *FluffDb) UpdateUserByName(username string, password string) error {
	rows, err := f.db.Query("UPDATE users SET password = ? WHERE username = ?", password, username)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

func (f *FluffDb) GetPgpByName(value string) (*Pgp, error) {
	var pgp Pgp
	// select packages
	rows, err := f.db.Query("SELECT key_name, private_key, public_key, passphrase FROM pgpkeys WHERE key_name = ?", value)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		if err = rows.Scan(&pgp.KeyName, &pgp.PrivateKey, &pgp.PublicKey, &pgp.PassPhrase); err != nil {
			return nil, fmt.Errorf("error reading pgp key: %v", err)
		}
	}
	return &pgp, nil
}

func (f *FluffDb) GetUserByName(value string) (*User, error) {
	var user User
	// select packages
	rows, err := f.db.Query("SELECT username, password, mode FROM users WHERE username = ?", value)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		if err = rows.Scan(&user.UserName, &user.Password, &user.Mode); err != nil {
			return nil, fmt.Errorf("error reading user name: %v", err)
		}
	}
	return &user, nil
}

func (f *FluffDb) DeletePgpByName(value string) error {
	// select packages
	rows, err := f.db.Query("DELETE FROM  pgpkeys WHERE key_name = ?", value)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

func (f *FluffDb) DeleteUserByName(value string) error {
	// select packages
	rows, err := f.db.Query("DELETE FROM  users WHERE username = ?", value)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

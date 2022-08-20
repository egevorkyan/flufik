package simpledb

import (
	"database/sql"
	"fmt"
	"github.com/egevorkyan/flufik/pkg/logging"
	_ "modernc.org/sqlite"
	"os"
)

type UserDb struct {
	db     *sql.DB
	dbpath string
	logger *logging.Logger
}

type PgpDb struct {
	db     *sql.DB
	dbpath string
	logger *logging.Logger
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

// CreateUserDb - Create internal Database
func CreateUserDb(path string, logger *logging.Logger) (*UserDb, error) {
	logger.Info("create user database")
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
	_, err = db.Exec(userCreateIndexes)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB tables: %v", err)
	}

	// create database indexes
	_, err = db.Exec(userCreateIndexes)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB indexes: %v", err)
	}

	return &UserDb{
		db:     db,
		dbpath: path,
		logger: logger,
	}, nil
}

// CreatePgpDb - Create internal Database
func CreatePgpDb(path string, logger *logging.Logger) (*PgpDb, error) {
	logger.Info("create pgp database")
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
	_, err = db.Exec(pgpCreateTables)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB tables: %v", err)
	}

	// create database indexes
	_, err = db.Exec(pgpCreateIndexes)
	if err != nil {
		return nil, fmt.Errorf("error creating Internal DB indexes: %v", err)
	}

	return &PgpDb{
		db:     db,
		dbpath: path,
		logger: logger,
	}, nil
}

// OpenUserDB opens a internal database SQLite database from file and return a
// pointer to the resulting struct.
func OpenUserDB(path string, logger *logging.Logger) (*UserDb, error) {
	logger.Info("open user database")
	// open database file
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can not open internal database: %v", err)
	}
	return &UserDb{
		db:     db,
		dbpath: path,
		logger: logger,
	}, nil
}

// OpenPgpDB opens a internal database SQLite database from file and return a
// pointer to the resulting struct.
func OpenPgpDB(path string, logger *logging.Logger) (*PgpDb, error) {
	logger.Info("open pgp database")
	// open database file
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("can not open internal database: %v", err)
	}
	return &PgpDb{
		db:     db,
		dbpath: path,
		logger: logger,
	}, nil
}

func (u *UserDb) Begin() (*sql.Tx, error) {
	return u.db.Begin()
}

func (u *UserDb) Close() error {
	if u.db != nil {
		return u.db.Close()
	}
	return nil
}

func (p *PgpDb) Begin() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *PgpDb) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PgpDb) InsertPgpKeys(keyName string, privateKey string, publicKey string, passPhrase string) error {
	p.logger.Info("insert pgp key to internal database")
	// insert pgp key
	statement, err := p.db.Prepare(sqlInsertPgpKey)
	if err != nil {
		return fmt.Errorf("can not prepare database statement: %v", err)
	}
	defer func(statement *sql.Stmt) {
		err = statement.Close()
		if err != nil {
			return
		}
	}(statement)

	_, err = statement.Exec(keyName, privateKey, publicKey, passPhrase)
	if err != nil {
		return fmt.Errorf("can not execute database statement: %v", err)
	}
	return nil
}

func (u *UserDb) InsertUsers(userName string, password string, mode string) error {
	u.logger.Info("insert user to internal database")
	// insert  user
	statement, err := u.db.Prepare(sqlInsertUser)
	if err != nil {
		return fmt.Errorf("can not prepare database statement: %v", err)
	}
	defer func(statement *sql.Stmt) {
		err = statement.Close()
		if err != nil {
			return
		}
	}(statement)

	_, err = statement.Exec(userName, password, mode)
	if err != nil {
		return fmt.Errorf("can not execute database statement: %v", err)
	}
	return nil
}

func (u *UserDb) UpdateUserByName(username string, password string) error {
	u.logger.Info("update user parameters")
	rows, err := u.db.Query("UPDATE users SET password = ? WHERE username = ?", password, username)
	if err != nil {
		return fmt.Errorf("user update failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

func (p *PgpDb) GetPgpByName(value string) (*Pgp, error) {
	p.logger.Info("requesting from database pgp key by name")
	var pgp Pgp
	// select packages
	rows, err := p.db.Query("SELECT key_name, private_key, public_key, passphrase FROM pgpkeys WHERE key_name = ?", value)
	if err != nil {
		return nil, fmt.Errorf("pgp key request failed: %v", err)
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
	fmt.Println(pgp)
	return &pgp, nil
}

func (u *UserDb) GetUserByName(value string) (*User, error) {
	u.logger.Info("requesting user from database by username")
	var user User
	// select packages
	rows, err := u.db.Query("SELECT username, password, mode FROM users WHERE username = ?", value)
	if err != nil {
		return nil, fmt.Errorf("user request failed: %v", err)
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

func (p *PgpDb) DeletePgpByName(value string) error {
	p.logger.Info("delete pgp key by name")
	// select packages
	rows, err := p.db.Query("DELETE FROM  pgpkeys WHERE key_name = ?", value)
	if err != nil {
		return fmt.Errorf("pgp key delete request failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

func (u *UserDb) DeleteUserByName(value string) error {
	u.logger.Info("delete user by username")
	// select packages
	rows, err := u.db.Query("DELETE FROM  users WHERE username = ?", value)
	if err != nil {
		return fmt.Errorf("user delete request failed: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)
	return nil
}

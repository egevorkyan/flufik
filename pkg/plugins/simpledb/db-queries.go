package simpledb

const (
	pgpCreateTables   = `CREATE TABLE IF NOT EXISTS pgpkeys ( key_name TEXT,  private_key TEXT, public_key TEXT, passphrase TEXT);`
	userCreateTables  = `CREATE TABLE IF NOT EXISTS users ( username TEXT,  password TEXT, mode TEXT);`
	pgpCreateIndexes  = `CREATE UNIQUE INDEX IF NOT EXISTS keynames ON pgpkeys (key_name);`
	userCreateIndexes = `CREATE UNIQUE INDEX IF NOT EXISTS usernames ON users (username);`
	sqlSelectPgpKey   = `SELECT
 key_name
 , private_key
 , public_key
 , passphrase
FROM pgpkeys;`
	sqlSelectUser = `SELECT
 username
 , password
 , mode
FROM users;`
	sqlInsertPgpKey = `INSERT INTO pgpkeys(key_name, private_key, public_key, passphrase) VALUES (?, ?, ?, ?);`
	sqlInsertUser   = `INSERT INTO users(username, password, mode) VALUES (?, ?, ?);`
)

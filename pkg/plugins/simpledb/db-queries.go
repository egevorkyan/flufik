package simpledb

const (
	sqlCreateTables = `CREATE TABLE IF NOT EXISTS pgpkeys ( key_name TEXT,  private_key TEXT, public_key TEXT, passphrase TEXT);
CREATE TABLE IF NOT EXISTS users ( username TEXT,  password TEXT, mode TEXT);`
	sqlCreateIndexes = `CREATE INDEX keynames ON pgpkeys (key_name);
CREATE INDEX usernames ON users (username);`
	sqlSelectPgpKey = `SELECT
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

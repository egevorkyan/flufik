package simpledb

const (
	flufikTableSqlSchema = `CREATE TABLE IF NOT EXISTS flufikTab (
		"idKey" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"key" TEXT,
		"privateKey" TEXT,
		"publicKey" TEXT,
		"token" TEXT
	  );`  //FlufikTable schema
	flufikTableSqlUnique = `CREATE UNIQUE index IF NOT EXISTS flufik_key on flufikTab(key)`
	DeleteFlufikSql      = `DELETE FROM flufikTab WHERE key = ?`
	appTableSqlSchema    = `CREATE TABLE IF NOT EXISTS appTab (
		"idKey" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"appname" TEXT,
        "appversion" TEXT,
		"apparch" TEXT,
		"apposversion" TEXT,
		"applocation" TEXT
	  );`  //AppTable schema
	DeleteAppSqlByVersion           = `DELETE FROM appTab WHERE ( appname = ? AND appversion = ? )`
	DeleteAppSqlByVersionByArch     = `DELETE FROM appTab WHERE ( appname = ? AND appversion = ? AND apparch = ? )`
	DeleteAppSqlByVersionByArchByOs = `DELETE FROM appTab WHERE ( appname = ? AND appversion = ? AND apparch = ? AND apposversion = ? )`
	DeleteAppSql                    = `DELETE FROM appTab WHERE ( appname = ? ) `

	internalFlufikInsertSchema = `INSERT OR ignore INTO flufikTab (key, privateKey, publicKey, token) VALUES (?, ?, ?, ?)`
	internalAppInsertSchema    = `INSERT OR ignore INTO appTab (appname, appversion, apparch, apposversion, applocation) VALUES (?, ?, ?, ?, ?)`
)

type Data struct {
	KeyValue        string
	PrivateKeyValue string
	PublicKeyValue  string
	TokenValue      string
}
type App struct {
	AppName      string
	AppVersion   string
	AppArch      string
	AppOsVersion string
	AppLocation  string
}

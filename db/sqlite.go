package db

import (
	"github.com/estenssoros/relay/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func connectSQLite(creds config.DBCreds) (*gorm.DB, error) {
	return gorm.Open("sqlite3", creds.Database)
}

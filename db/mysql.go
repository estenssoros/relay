package db

import (
	"fmt"

	"github.com/estenssoros/relay/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func connectMySQL(creds config.DBCreds) (*gorm.DB, error) {
	connectionURL := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", creds.User, creds.Password, creds.Host, creds.Database)
	return gorm.Open("mysql", connectionURL)
}

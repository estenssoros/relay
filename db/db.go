package db

import (
	"github.com/estenssoros/goflow/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

var Connection *gorm.DB

func init() {
	conn, err := Connect()
	if err != nil {
		logrus.Fatal(err)
	}
	Connection = conn
}

func Connect() (*gorm.DB, error) {
	sqlConn := config.DefaultConfig.Core.SQLConn
	db, err := gorm.Open("sqlite3", sqlConn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

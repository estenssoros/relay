package db

import (
	"github.com/estenssoros/goflow/config"
	"github.com/jinzhu/gorm"
)

func Connect() (*gorm.DB, error) {
	sqlConn := config.DefaultConfig.Core.SQLConn
	db, err := gorm.Open("sqlite3", sqlConn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

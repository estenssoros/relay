package db

import (
	"github.com/estenssoros/relay/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

var Connection *gorm.DB

func init() {
	conn, err := Connect()
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "connect"))
	}
	Connection = conn
}

func Connect() (*gorm.DB, error) {
	switch config.DefaultConfig.DBCreds.Flavor {
	case "sqlite":
		return connectSQLite(config.DefaultConfig.DBCreds)
	case "mysql":
		return connectMySQL(config.DefaultConfig.DBCreds)
	case "postrgres":
		return connectPostgres(config.DefaultConfig.DBCreds)
	default:
		return nil, errors.Errorf("not supported: %s", config.DefaultConfig.DBCreds.Flavor)
	}
}

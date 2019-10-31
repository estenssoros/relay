package db

import (
	"fmt"

	"github.com/estenssoros/relay/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func connectPostgres(creds config.DBCreds) (*gorm.DB, error) {
	connectionURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		creds.Host,
		creds.Port,
		creds.User,
		creds.Password,
		creds.Database,
	)
	return gorm.Open("postgres", connectionURL)
}

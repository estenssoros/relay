package models

import (
	"github.com/estenssoros/goflow/connection"
)

type Connection struct {
	Base
	ConnName string
	ConnType connection.Type
	Host     string
	Schema   string
	Login    string
	Password string
	Port     int
}

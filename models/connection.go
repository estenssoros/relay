package models

// Connection placeholder to store information about different database instances
// connection information. The idea here is that scripts use references to
// database instances (conn_id) instead of hard coding hostname, logins and
// passwords when using operators or hooks.
type Connection struct {
	ID       int
	ConnID   string
	ConnType string
	Host     string
	Schema   string
	Login    string
	Password string
	Port     int
	Extra    string
}

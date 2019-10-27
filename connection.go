package goflow

var (
	// MySQLType mysql connection type
	MySQLType ConnectionType = "mysql"
	// PostgresType postgres connection type
	PostgresType ConnectionType = "postgres"
	// SnowFlakeType snowflake connection type
	SnowFlakeType ConnectionType = "snowflake"
	// MsSQLType mysql connetion type
	MsSQLType ConnectionType = "mssql"
	// ODBCType odbc connection type
	ODBCType ConnectionType = "odbc"
	// OracleType oracle connection type
	OracleType ConnectionType = "orcale"
	// SQLlite sqlite connection type
	SQLlite ConnectionType = "sqlite"
)
var (
	// S3Type s3 connection type
	S3Type ConnectionType = "s3"
	// FTPType ftp connection type
	FTPType ConnectionType = "ftp"
)

// ConnectionType type of connection
type ConnectionType string

// GetConnection gets a connection
func GetConnection(connectionID string) (ConnectionInterface, error) {
	return nil, nil
}

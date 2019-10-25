package connection

var (
	MySQLType     Type = "mysql"
	PostgresType  Type = "postgres"
	SnowFlakeType Type = "snowflake"
	MsSQL         Type = "mssql"
	ODBCType      Type = "odbc"
	OracleType    Type = "orcale"
	SQLlite       Type = "sqlite"
)
var (
	S3Type  Type = "s3"
	FTPType Type = "ftp"
)

type Type string

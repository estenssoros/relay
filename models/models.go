package models

// Migrations default models to migrate to database
var Migrations = []interface{}{
	&Connection{},
	&DAG{},
	&DagRun{},
	&TaskInstance{},
}

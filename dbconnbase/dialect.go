package dbconnbase

import "fmt"

type DialectType uint8

const (
	NENHUM     DialectType = 0
	FIREBIRD               = 1
	INTERBASE              = 2
	MYSQL                  = 3
	ORACLE                 = 4
	POSTGRESQL             = 5
	SQLSERVER              = 6
	SQLITE                 = 7
)

var DialectList = [7]DialectType{FIREBIRD, INTERBASE, MYSQL, ORACLE, POSTGRESQL, SQLSERVER, SQLITE}
var DialectName = [7]string{"Firebird", "Interbase", "MySQL", "Oracle", "PostgreSQL", "SQL Server", "SQLite"}
var DialectDrive = [7]string{"firebird", "interbase", "mysql", "oracle", "pgx", "sqlserver", "sqlite"}
var DialectLow = [7]string{"FIRE", "INTER", "MYSQL", "ORA", "PG", "SQLSERVER", "SQLITE"}

func (d DialectType) String() string {
	switch d {
	case FIREBIRD:
		return "firebird"
	case INTERBASE:
		return "interbase"
	case MYSQL:
		return "mysql"
	case ORACLE:
		return "oracle"
	case POSTGRESQL:
		return "pgx"
	case SQLSERVER:
		return "sqlserver"
	case SQLITE:
		return "sqlite"
	default:
		return fmt.Sprintf("%d", int(d))
	}
}
func DialectDriveFromString(s string) DialectType {
	for i, value := range DialectDrive {
		if value == s {
			return DialectList[i]
		}
	}
	return 0
}
func DialectLowFromString(s string) DialectType {
	for i, value := range DialectLow {
		if value == s {
			return DialectList[i]
		}
	}
	return 0
}
func DialectNameFromString(s string) DialectType {
	for i, value := range DialectName {
		if value == s {
			return DialectList[i]
		}
	}
	return 0
}

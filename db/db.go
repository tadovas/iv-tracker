package db

import (
	"database/sql"
	"flag"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Flags struct {
	Host     string
	Username string
	Password string
	Name     string
}

func RegisterFlags(dbFlags *Flags) {
	flag.StringVar(&dbFlags.Host, "db.host", "localhost", "Host name of database")
	flag.StringVar(&dbFlags.Username, "db.username", "", "Username for database")
	flag.StringVar(&dbFlags.Password, "db.password", "", "Password for database")
	flag.StringVar(&dbFlags.Name, "db.name", "", "Database name")
}

func Setup(dbFlags Flags) (*sql.DB, error) {
	config := mysql.NewConfig()
	config.Net = "tcp"
	config.Addr = dbFlags.Host
	config.User = dbFlags.Username
	config.Passwd = dbFlags.Password
	config.DBName = dbFlags.Name
	config.ParseTime = true

	return sql.Open("mysql", config.FormatDSN())
}

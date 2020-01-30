package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/tadovas/iv-tracker/log"

	"github.com/pressly/goose"
)

func init() {
	goose.SetLogger(loggerAdapter{})
}

func Migrate(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("mysql"); err != nil {
		return fmt.Errorf("goose set dialect error: %w", err)
	}
	if err := goose.Status(db, migrationsDir); err != nil {
		return fmt.Errorf("goose status error: %w", err)
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}
	return nil
}

type loggerAdapter struct {
}

func (_ loggerAdapter) Fatal(v ...interface{}) {
	log.Error(v...)
	os.Exit(2)
}

func (_ loggerAdapter) Fatalf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(2)
}

func (_ loggerAdapter) Print(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	log.Infof(format, v...)
}

func (_ loggerAdapter) Println(v ...interface{}) {
	log.Info(v...)
}

func (_ loggerAdapter) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}

var _ goose.Logger = loggerAdapter{}

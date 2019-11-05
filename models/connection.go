package models

import (
	"errors"
	"fmt"

	"github.com/MelchiSalins/go-auth/pkg/app"
	"github.com/jinzhu/gorm"

	// Required by GORM to connect to postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// Required by GORM to connect to sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// ErrorInvalidDBName returned when an invalid database name is provided
	ErrorInvalidDBName = errors.New("models: Invalid database name provided")
)

// OpenDBConn Returns sqlite GORM DB connection object
// Accepted values for dbname is sqlite or postgres
func OpenDBConn(dbname string, logmode bool) (*gorm.DB, error) {
	switch dbname {
	case "sqlite":
		db, err := gorm.Open("sqlite3", "user.sqlite")
		db.LogMode(logmode)
		return db, err
	case "postgres", "pg":
		c := fmt.Sprintf("host=%s port=%d sslmode=%s user=%s password=%s", app.DBHost, app.DBPort, app.DBSSLMode, app.DBUser, app.DBPass)
		db, err := gorm.Open("postgres", c)
		db.LogMode(logmode)
		return db, err
	default:
		return nil, ErrorInvalidDBName
	}
}

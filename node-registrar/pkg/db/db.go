package db

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host        string
	Port        int
	DBName      string
	User        string
	Password    string
	SSLMode     string
	SqlLogLevel logger.LogLevel
	MaxConns    int
}

// PostgresDatabase postgres db client
type DataBase struct {
	gormDB     *gorm.DB
	connString string
}

var ErrRecordNotFound = errors.New("could not find any records")

func NewDB(c Config) (db DataBase, err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(c.SqlLogLevel),
	})
	if err != nil {
		return db, errors.Wrapf(err, "Failed to connect to the database: %v", err)
	}

	sql, err := gormDB.DB()
	if err != nil {
		return db, errors.Wrap(err, "failed to configure DB connection")
	}

	sql.SetMaxIdleConns(3)
	sql.SetMaxOpenConns(c.MaxConns)

	db = DataBase{gormDB, dsn}

	if err := db.gormDB.AutoMigrate(
		&Farm{},
		&Node{},
	); err != nil {
		return db, errors.Wrap(err, "failed to migrate tables")
	}
	return
}

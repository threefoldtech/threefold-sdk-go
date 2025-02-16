package db

import (
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/pkg/errors"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/metrics"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	PostgresHost     string
	PostgresPort     uint64
	DBName           string
	PostgresUser     string
	PostgresPassword string
	SSLMode          string
	SqlLogLevel      logger.LogLevel
	MaxOpenConns     uint64
	MaxIdleConns     uint64
}

// PostgresDatabase postgres db client
type Database struct {
	gormDB     *gorm.DB
	connString string
	metrics    *metrics.Metrics
}

var (
	ErrRecordNotFound      = errors.New("could not find any records")
	ErrRecordAlreadyExists = errors.New("record already exists")
)

func NewDB(c Config, m *metrics.Metrics) (Database, error) {
	db, err := openDatabase(c, m)
	if err != nil {
		return Database{}, err
	}

	sql, err := db.gormDB.DB()
	if err != nil {
		return Database{}, errors.Wrap(err, "failed to configure DB connection")
	}

	sql.SetMaxIdleConns(int(c.MaxIdleConns))
	sql.SetMaxOpenConns(int(c.MaxOpenConns))

	err = db.autoMigrate()
	if err != nil {
		return Database{}, err
	}

	return db, sql.Ping()
}

func openDatabase(c Config, m *metrics.Metrics) (db Database, err error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.DBName, c.SSLMode)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(c.SqlLogLevel),
	})
	if err != nil {
		return db, errors.Wrapf(err, "Failed to connect to the database: %v", err)
	}

	return Database{gormDB, dsn, m}, nil
}

func (db Database) autoMigrate() error {
	if err := db.gormDB.AutoMigrate(
		&Account{},
		&Farm{},
		&Node{},
		&UptimeReport{},
		&ZosVersion{},
	); err != nil {
		return errors.Wrap(err, "failed to migrate tables")
	}
	return nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.PostgresHost) == "" {
		return errors.New("invalid postgres host, postgres host should not be empty")
	}

	if net.ParseIP(c.PostgresHost) == nil {
		if _, err := net.LookupHost(c.PostgresHost); err == nil {
			return errors.Wrapf(err, "invalid postgres host %s, failed to parse or lookup host", c.PostgresHost)
		}
	}

	if c.PostgresPort < 1 && c.PostgresPort > 65535 {
		return errors.Errorf("invalid postgres port %d, postgres port should be in the valid port range 1–65535", c.PostgresPort)
	}

	if strings.TrimSpace(c.DBName) == "" {
		return errors.New("invalid database name, database name should not be empty")
	}

	if strings.TrimSpace(c.PostgresUser) == "" {
		return errors.New("invalid postgres user, postgres user should not be empty")
	}

	if strings.TrimSpace(c.PostgresPassword) == "" {
		return errors.New("invalid postgres password, postgres password should not be empty")
	}

	sslModes := []string{"disable", "require", "verify-ca", "verify-full"}
	if !slices.Contains(sslModes, c.SSLMode) {
		return errors.New(fmt.Sprintf("invalid ssl mode %s, ssl mode should be one of %v", c.SSLMode, sslModes))
	}

	sqlLogLevel := map[int]string{1: "Silent", 2: "Error", 3: "Warn", 4: "Info"}
	if c.SqlLogLevel < 0 || c.SqlLogLevel > 4 {
		return errors.New(fmt.Sprintf("invalid sql log level %d, sql log level should be one of %v", c.SqlLogLevel, sqlLogLevel))
	}

	return nil
}

func (db Database) Close() error {
	sql, err := db.gormDB.DB()
	if err != nil {
		return err
	}
	if err := sql.Close(); err != nil {
		return err
	}
	return nil
}

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/metrics"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/server"
	"gorm.io/gorm/logger"
)

type flags struct {
	db.Config
	debug       bool
	version     bool
	domain      string
	serverPort  uint
	network     string
	adminTwinID uint64
}

var (
	commit  string
	version string
)

func main() {
	if err := Run(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func Run() error {
	f := flags{}
	var sqlLogLevel int
	flag.StringVar(&f.PostgresHost, "postgres-host", "", "postgres host")
	flag.Uint64Var(&f.Config.PostgresPort, "postgres-port", 5432, "postgres port")
	flag.StringVar(&f.DBName, "postgres-db", "", "postgres database")
	flag.StringVar(&f.PostgresUser, "postgres-user", "", "postgres username")
	flag.StringVar(&f.PostgresPassword, "postgres-password", "", "postgres password")
	flag.StringVar(&f.SSLMode, "ssl-mode", "disable", "postgres ssl mode[disable, require, verify-ca, verify-full]")
	flag.IntVar(&sqlLogLevel, "sql-log-level", 2, "sql logger level")
	flag.Uint64Var(&f.MaxOpenConns, "max-open-conn", 3, "max open sql connections")
	flag.Uint64Var(&f.MaxIdleConns, "max-idle-conn", 3, "max idle sql connections")

	flag.BoolVar(&f.version, "v", false, "shows the package version")
	flag.BoolVar(&f.debug, "debug", false, "allow debug logs")
	flag.UintVar(&f.serverPort, "server-port", 8080, "server port")
	flag.StringVar(&f.domain, "domain", "", "domain on which the server will be served")
	flag.StringVar(&f.network, "network", "dev", "the registrar network")
	flag.Uint64Var(&f.adminTwinID, "admin-twin-id", 0, "admin twin ID")

	flag.Parse()
	f.SqlLogLevel = logger.LogLevel(sqlLogLevel)

	if f.version {
		log.Info().Str("version", version).Str("commit", commit).Send()
		return nil
	}

	if err := f.validate(); err != nil {
		return err
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if f.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	m := metrics.NewMetrics()
	if err := m.Register(); err != nil {
		return errors.Wrap(err, "failed to register metrics")
	}

	db, err := db.NewDB(f.Config, m)
	if err != nil {
		return errors.Wrap(err, "failed to open database with the specified configurations")
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()

	s, err := server.NewServer(db, f.network, f.adminTwinID, m)
	if err != nil {
		return errors.Wrap(err, "failed to start gin server")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Info().Msg("server is running on port :8080")

	err = s.Run(quit, fmt.Sprintf("%s:%d", f.domain, f.serverPort))
	if err != nil {
		return errors.Wrap(err, "failed to run gin server")
	}

	return nil
}

func (f flags) validate() error {
	if f.serverPort < 1 || f.serverPort > 65535 {
		return errors.Errorf("invalid port %d, server port should be in the valid port range 1–65535", f.serverPort)
	}

	if strings.TrimSpace(f.domain) == "" {
		return errors.New("invalid domain name, domain name should not be empty")
	}
	if _, err := net.LookupHost(f.domain); err != nil {
		return errors.Wrapf(err, "invalid domain %s", f.domain)
	}

	return f.Config.Validate()
}

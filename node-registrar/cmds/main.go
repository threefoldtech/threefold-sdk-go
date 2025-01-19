package cmds

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/server"
	"gorm.io/gorm/logger"
)

type flags struct {
	db.Config
	debug      bool
	version    bool
	domain     string
	serverPort uint
}

var (
	commit  string
	version string
)

func Run() {
	f := flags{}
	var sqlLogLevel int
	flag.StringVar(&f.PostgresHost, "postgres-host", "", "postgres host")
	flag.UintVar(&f.Config.PostgresPort, "postgres-port", 5432, "postgres port")
	flag.StringVar(&f.DBName, "postgres-db", "", "postgres database")
	flag.StringVar(&f.PostgresUser, "postgres-user", "", "postgres username")
	flag.StringVar(&f.PostgresPassword, "postgres-password", "", "postgres password")
	flag.StringVar(&f.SSLMode, "ssl-mode", "disable", "postgres ssl mode[disable, require, verify-ca, verify-full]")
	flag.IntVar(&sqlLogLevel, "sql-log-level", 2, "sql logger level")
	flag.UintVar(&f.MaxConns, "max-conn", 3, "max sql connections")

	flag.BoolVar(&f.version, "v", false, "shows the package version")
	flag.BoolVar(&f.debug, "debug", false, "allow debug logs")
	flag.UintVar(&f.serverPort, "server-port", 8080, "server port")
	flag.StringVar(&f.domain, "domain", "", "domain on which the server will be served")

	flag.Parse()
	f.SqlLogLevel = logger.LogLevel(sqlLogLevel)

	if f.version {
		log.Info().Str("version", version).Str("commit", commit).Send()
		return
	}

	if err := f.validate(); err != nil {
		log.Error().Err(err).Send()
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if f.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	db, err := db.NewDB(f.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database with the specified configurations")
	}

	s, err := server.NewServer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start gin server")
	}

	err = s.Run(fmt.Sprintf("%s:%d", f.domain, f.serverPort))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run gin server")
	}
}

func (f flags) validate() error {
	if f.serverPort < 1 && f.serverPort > 65535 {
		return errors.Errorf("invalid port %d, server port should be in the valid port range 1–65535", f.serverPort)
	}

	if strings.TrimSpace(f.domain) == "" {
		return errors.New("invalid domain name, domain name should not be empty")
	}
	if _, err := net.LookupHost(f.domain); err != nil {
		return errors.Wrapf(err, "invalid domain %s", f.PostgresHost)
	}

	return f.Config.Validate()
}

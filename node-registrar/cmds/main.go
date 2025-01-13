package cmds

import (
	"errors"
	"flag"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/server"
	"gorm.io/gorm/logger"
)

type flags struct {
	db.Config
	debug   bool
	version bool
	domain  string
	port    int
}

var (
	commit  string
	version string
)

func Run() {
	f := flags{}
	var sqlLogLevel int
	flag.StringVar(&f.Host, "postgres-host", "", "postgres host")
	flag.IntVar(&f.Config.Port, "postgres-port", 5432, "postgres port")
	flag.StringVar(&f.DBName, "postgres-db", "", "postgres database")
	flag.StringVar(&f.User, "postgres-user", "", "postgres username")
	flag.StringVar(&f.Password, "postgres-password", "", "postgres password")
	flag.StringVar(&f.SSLMode, "ssl-mode", "disable", "postgres ssl mode[disable, require, verify-ca, verify-full]")
	flag.IntVar(&sqlLogLevel, "sql-log-level", 2, "sql logger level")

	flag.BoolVar(&f.version, "v", false, "shows the package version")
	flag.BoolVar(&f.debug, "debug", false, "allow debug logs")
	flag.IntVar(&f.port, "server-port", 443, "server port")
	flag.StringVar(&f.domain, "domain", "", "domain on which the server will be served")

	flag.Parse()
	f.SqlLogLevel = logger.LogLevel(sqlLogLevel)

	if f.version {
		log.Info().Str("version", version).Str("commit", commit).Send()
		return
	}
	if f.domain == "" {
		log.Fatal().Err(errors.New("domain is required"))
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if f.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	db, err := db.NewDB(f.Config)
	if err != nil {
		log.Fatal().Msg("failed to open database with the specified configurations")
	}

	s, err := server.NewServer(db)
	if err != nil {
		log.Fatal().Msg("failed to start gin server")
	}

	err = s.Router.Run(fmt.Sprintf("%s:%d", f.domain, f.port))
	if err != nil {
		log.Fatal().Msg("failed to run gin server")
	}
}

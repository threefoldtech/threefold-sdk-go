package server

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
)

type Server struct {
	router      *gin.Engine
	db          db.Database
	network     string
	adminTwinID uint64
}

func NewServer(db db.Database, network string, adminTwinID uint64) (s Server, err error) {
	router := gin.Default()

	s = Server{router, db, network, adminTwinID}
	s.SetupRoutes()

	return
}

func (s Server) Run(quit chan os.Signal, addr string) error {
	server := &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-quit

		log.Info().Msg("server is shutting down")
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shut down server gracefully")
		}
	}()

	err := server.ListenAndServe()
	wg.Wait()

	return err
}

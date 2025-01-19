package server

import (
	"github.com/gin-gonic/gin"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
)

type Server struct {
	router *gin.Engine
	db     db.DataBase
}

func NewServer(db db.DataBase) (s Server, err error) {
	router := gin.Default()

	s = Server{router, db}
	s.SetupRoutes()

	return
}

func (s Server) Run(addr string) error {
	return s.router.Run(addr)
}

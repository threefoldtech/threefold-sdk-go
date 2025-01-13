package server

import (
	"github.com/gin-gonic/gin"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
)

type Server struct {
	Router *gin.Engine
	DB     db.DataBase
}

func NewServer(db db.DataBase) (s Server, err error) {
	router := gin.Default()

	s = Server{router, db}
	s.SetupRoutes()

	return
}

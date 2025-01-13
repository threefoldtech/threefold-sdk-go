package server

import (
	_ "github.com/threefoldtech/tfgrid-sdk-go/node-registrar/docs"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func (s *Server) SetupRoutes() {
	s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	farmRoutes := s.Router.Group("farms")
	farmRoutes.GET("/", s.listFarmsHandler)
	farmRoutes.GET("/:farm_id", s.getFarmHandler)
	farmRoutes.POST("/", s.createFarmHandler)
	farmRoutes.PATCH("/", s.updateFarmsHandler)

	nodeRoutes := s.Router.Group("nodes")
	nodeRoutes.GET("/", s.listNodesHandler)
	nodeRoutes.GET("/:node_id", s.getNodeHandler)
	nodeRoutes.POST("/", s.registerNodeHandler)
	nodeRoutes.POST("/:node_id/uptime", s.uptimeHandler)
	nodeRoutes.POST("/:node_id/consumption", s.storeConsumptionHandler)
}

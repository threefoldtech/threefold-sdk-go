package server

import (
	_ "github.com/threefoldtech/tfgrid-sdk-go/node-registrar/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) SetupRoutes() {
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := s.router.Group("v1")

	// farms routes
	farmRoutes := v1.Group("farms")
	farmRoutes.GET("/", s.listFarmsHandler)
	farmRoutes.GET("/:farm_id", s.getFarmHandler)
	farmRoutes.POST("/", s.createFarmHandler)

	farmRoutes.Use(s.createAuthFarmMiddleware(s.network))
	farmRoutes.PATCH("/:farm_id", s.updateFarmsHandler)

	// nodes routes
	nodeRoutes := v1.Group("nodes")
	nodeRoutes.GET("/", s.listNodesHandler)
	nodeRoutes.GET("/:node_id", s.getNodeHandler)
	nodeRoutes.POST("/", s.registerNodeHandler)

	nodeRoutes.Use(s.createAuthNodeMiddleware(s.network))
	nodeRoutes.POST("/:node_id/uptime", s.uptimeReportHandler)
	nodeRoutes.POST("/:node_id/consumption", s.storeConsumptionHandler)

	// Account routes
	accountRoutes := v1.Group("accounts")
	accountRoutes.POST("/", s.createAccountHandler)
	accountRoutes.GET("/:twin_id", s.getAccountHandler)
	accountRoutes.PATCH("/:twin_id", s.updateAccountHandler)

}

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
	// protected by farmer key
	farmRoutes.Use(s.AuthMiddleware())
	farmRoutes.POST("/", s.createFarmHandler)
	farmRoutes.PATCH("/:farm_id", s.updateFarmsHandler)

	// nodes routes
	nodeRoutes := v1.Group("nodes")
	nodeRoutes.GET("/", s.listNodesHandler)
	nodeRoutes.GET("/:node_id", s.getNodeHandler)
	// protected by node key
	nodeRoutes.Use(s.AuthMiddleware())
	nodeRoutes.POST("/", s.registerNodeHandler)
	nodeRoutes.PATCH("/:node_id", s.updateNodeHandler)
	nodeRoutes.POST("/:node_id/uptime", s.uptimeReportHandler)

	// Account routes
	accountRoutes := v1.Group("accounts")
	accountRoutes.POST("/", s.createAccountHandler)
	accountRoutes.GET("/", s.getAccountHandler)
	// protected by farmer key
	accountRoutes.Use(s.AuthMiddleware())
	accountRoutes.PATCH("/:twin_id", s.updateAccountHandler)
}

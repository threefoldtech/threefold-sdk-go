package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/subi"
)

var (
	DevNetwork  = "dev"
	QaNetwork   = "qa"
	TestNetwork = "test"
	MainNetwork = "main"

	SubstrateURLs = map[string][]string{
		DevNetwork: {
			"wss://tfchain.dev.grid.tf/ws",
			"wss://tfchain.dev.grid.tf:443",
			"wss://tfchain.02.dev.grid.tf/ws",
			"wss://tfchain.02.dev.grid.tf:443",
		},
		QaNetwork: {
			"wss://tfchain.qa.grid.tf/ws",
			"wss://tfchain.qa.grid.tf:443",
			"wss://tfchain.02.qa.grid.tf/ws",
			"wss://tfchain.02.qa.grid.tf:443",
		},
		TestNetwork: {
			"wss://tfchain.test.grid.tf/ws",
			"wss://tfchain.test.grid.tf:443",
			"wss://tfchain.02.test.grid.tf/ws",
			"wss://tfchain.02.test.grid.tf:443",
		},
		MainNetwork: {
			"wss://tfchain.grid.tf/ws",
			"wss://tfchain.grid.tf:443",
			"wss://tfchain.02.grid.tf/ws",
			"wss://tfchain.02.grid.tf:443",
		},
	}
)

func (server Server) createAuthFarmMiddleware(network string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		pubKey, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			fmt.Println("Error decoding token:", err)
			return
		}
		manager := subi.NewManager(SubstrateURLs[network]...)

		sub, err := manager.SubstrateExt()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "could not get substrate client"})
			return
		}

		twinID, err := sub.GetTwinByPubKey(pubKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to get twin id with public key"})
			return
		}

		farmID := c.Param("farm_id")

		id, err := strconv.ParseUint(farmID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid farm_id: %v", err.Error())})
			return
		}

		farm, err := server.db.GetFarm(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid farm_id: %v", err.Error())})
			return
		}

		if twinID != uint32(farm.TwinID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or public key"})
			return
		}

		c.Next()
	}
}

func (server Server) createAuthNodeMiddleware(network string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		pubKey, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		manager := subi.NewManager(SubstrateURLs[network]...)
		sub, err := manager.SubstrateExt()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		twinID, err := sub.GetTwinByPubKey(pubKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		nodeID := c.Param("node_id")
		id, err := strconv.ParseUint(nodeID, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		node, err := server.db.GetNode(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		if twinID != uint32(node.TwinID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid twin id or public key"})
			return
		}

		c.Next()
	}
}

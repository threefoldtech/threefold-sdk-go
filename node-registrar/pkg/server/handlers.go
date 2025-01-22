package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
)

// @Summary list farms
// @Description list farms with specific filter
// @Accept  json
// @Produce  json
// @Param farm_name query string false "farm name"
// @Param farm_id query uint64 false "farm id"
// @Param twin_id query uint64 false "twin id"
// @Param page query int false "Page number"
// @Param size query int false "Max result per page"
// @Success 200 {object} []db.Farm
// @Failure 400 {object} error
// @Router /farms/ [get]
func (s Server) listFarmsHandler(c *gin.Context) {
	var filter db.FarmFilter
	limit := db.DefaultLimit()

	err := parseQueryParams(c, &limit, &filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	farms, err := s.db.ListFarms(filter, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"farms": farms,
	})
}

// @Summary get farm
// @Description get a farm with specific id
// @Accept  json
// @Produce  json
// @param farm_id path uint64 true "farm id"
// @Success 200 {object} db.Farm
// @Failure 404 {object} db.ErrRecordNotFound
// @Failure 400 {object} error
// @Router /farm/{farm_id} [get]
func (s Server) getFarmHandler(c *gin.Context) {
	farmID := c.Param("farm_id")

	id, err := strconv.ParseUint(farmID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid farm_id: %v", err.Error())})
		return
	}

	farm, err := s.db.GetFarm(id)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"farm": farm,
	})
}

// @Summary create a farm
// @Description creates a farm
// @Accept  json
// @Produce  json
// @Param farm_id body uint64 true "farm id"
// @Param farm_name body uint64 true "farm name"
// @Param twin_id body uint64 true "twin id"
// @Param dedicated body bool false "dedicated farm"
// @Param farm_free_ips body uint64 false "farm free ips"
// @Success 200 {object} db.Farm
// @Failure 400 {object} error
// @Failure 409 {object} db.ErrRecordAlreadyExists
// @Router /farms/ [post]
func (s Server) createFarmHandler(c *gin.Context) {
	var farm db.Farm

	if err := c.ShouldBindJSON(&farm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse farm info: %v", err.Error())})
		return
	}

	err := s.db.CreateFarm(farm)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordAlreadyExists) {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Farm created successfully",
		"farm":    farm,
	})
}

// @Summary update a farm
// @Description update a farm
// @Accept  json
// @Produce  json
// @Param farm_id body uint64 true "farm id"
// @Param farm_name body uint64 false "farm name"
// @Param twin_id body uint64 false "twin id"
// @Param dedicated body bool false "dedicated farm"
// @Param farm_free_ips body uint64 false "farm free ips"
// @Success 200 {object} db.Farm
// @Failure 400 {object} error
// @Failure 404 {object} db.ErrRecordNotFound
// @Router /farms/ [patch]
func (s Server) updateFarmsHandler(c *gin.Context) {
	var farm db.Farm
	farmID := c.Param("farm_id")

	id, err := strconv.ParseUint(farmID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid farm_id: %v", err.Error())})
		return
	}

	if err := c.ShouldBindJSON(&farm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse farm info: %v", err.Error())})
		return
	}

	if farm.FarmID != 0 && farm.FarmID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "farm id does not match farm id in the request"})
		return
	}

	err = s.db.UpdateFarm(id, farm)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Farm was updated successfully",
	})
}

// @Summary list nodes
// @Description list nodes with specific filter
// @Accept  json
// @Produce  json
// @Param node_id query uint64 false "node id"
// @Param farm_id query uint64 false "farm id"
// @Param twin_id query uint64 false "twin id"
// @Param status query string false "node status"
// @Param healthy query bool false "is node healthy"
// @Param page query int false "Page number"
// @Param size query int false "Max result per page"
// @Success 200 {object} []db.Node
// @Failure 400 {object} error
// @Router /nodes/ [get]
func (s Server) listNodesHandler(c *gin.Context) {
	var filter db.NodeFilter
	limit := db.DefaultLimit()

	err := parseQueryParams(c, &limit, &filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(filter, limit)

	nodes, err := s.db.ListNodes(filter, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

// @Summary get node
// @Description get a node with specific id
// @Accept  json
// @Produce  json
// @param node_id path uint64 false "node id"
// @Success 200 {object} db.Node
// @Failure 404 {object} db.ErrRecordNotFound
// @Failure 400 {object} error
// @Router /nodes/{node_id} [get]
func (s Server) getNodeHandler(c *gin.Context) {
	nodeID := c.Param("node_id")

	id, err := strconv.ParseUint(nodeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node id"})
		return
	}

	node, err := s.db.GetNode(id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"node": node,
	})
}

// @Summary register a node
// @Description register a node
// @Accept  json
// @Produce  json
// @Param node_id body uint64 true "node id"
// @Param farm_id body uint64 true "farm id"
// @Param twin_id body uint64 true "twin id"
// @Param features body []string true "node features "
// @Param status body string false "node status"
// @Param healthy body bool false "node healthy"
// @Param dedicated body bool false "node dedicated"
// @Param rented body bool false "node rented"
// @Param rentable body bool false "node rentable"
// @Param price_usd body float64 false "price in usd"
// @Param uptime body db.Uptime false "uptime report"
// @Param consumption body db.Consumption false "consumption report"
// @Success 200 {object} db.Node
// @Failure 400 {object} error
// @Failure 409 {object} db.ErrRecordAlreadyExists
// @Router /nodes/ [post]
func (s Server) registerNodeHandler(c *gin.Context) {
	var node db.Node

	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.db.RegisterNode(node)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordAlreadyExists) {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "node registered successfully",
		"node":    node,
	})
}

// @Summary uptime report
// @Description save uptime report of a node
// @Accept  json
// @Produce  json
// @Param node_id query uint64 true "node id"
// @Param uptime body db.Uptime false "uptime report"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 404 {object} db.ErrRecordNotFound
// @Router /nodes/{node_id}/uptime [post]
func (s Server) uptimeHandler(c *gin.Context) {
	nodeID := c.Param("node_id")

	id, err := strconv.ParseUint(nodeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node id"})
		return
	}

	var report struct {
		Uptime db.Uptime `json:"uptime"`
	}

	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.db.Uptime(id, report.Uptime)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "uptime report received successfully",
	})
}

// @Summary consumption report
// @Description save consumption report of a node
// @Accept  json
// @Produce  json
// @Param node_id query uint64 true "node id"
// @Param consumption body db.Consumption false "consumption report"
// @Success 200 {object} string
// @Failure 400 {object} error
// @Failure 404 {object} db.ErrRecordNotFound
// @Router /nodes/{node_id}/uptime [post]
func (s Server) storeConsumptionHandler(c *gin.Context) {
	nodeID := c.Param("node_id")

	id, err := strconv.ParseUint(nodeID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node id"})
		return
	}

	var consumption db.Consumption

	if err := c.ShouldBindJSON(&consumption); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.db.Consumption(id, consumption)
	if err != nil {
		status := http.StatusBadRequest

		if errors.Is(err, db.ErrRecordNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "consumption report received successfully",
	})
}

func parseQueryParams(c *gin.Context, types_ ...interface{}) error {
	for _, type_ := range types_ {
		if err := c.ShouldBindQuery(type_); err != nil {
			return fmt.Errorf("failed to bind query params to %T: %w", type_, err)
		}
	}
	return nil
}

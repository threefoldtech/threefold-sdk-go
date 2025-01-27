package server

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

type UptimeReportRequest struct {
	NodeID    uint64        `json:"node_id" binding:"required"`
	Uptime    time.Duration `json:"uptime" binding:"required"`
	Timestamp time.Time     `json:"timestamp" binding:"required"`
}

// @Summary uptime report
// @Description save uptime report of a node
// @Accept  json
// @Produce  json
// @Param node_id path uint64 true "node id"
// @Param request body UptimeReportRequest true "uptime report request"
// @Success 201 {object} map[string]string "message: uptime reported successfully"
// @Failure 400 {object} map[string]string "error: error message"
// @Failure 404 {object} map[string]string "error: node not found"
// @Router /nodes/{node_id}/uptime [post]
func (s *Server) uptimeReportHandler(c *gin.Context) {
	var req UptimeReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get node and last report
	_, err := s.db.GetNode(req.NodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	// Detect restarts
	// Validate report timing (40min ± 5min window)
	// Maybe aggregate reports here and store total uptime?
	// The total uptime should accumulate unless the node restarts, which is detected when the reported uptime is less than the previous value.

	// Create report record
	report := &db.UptimeReport{
		NodeID:    req.NodeID,
		Duration:  req.Uptime,
		Timestamp: req.Timestamp,
	}

	err = s.db.CreateUptimeReport(report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save report"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "uptime reported successfully"})
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

// AccountRequest represents the request body for account operations
type AccountCreationRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}

const (
	MaxTimestampDelta = 2 * time.Second
)

// Challenge is uniquely tied to both the timestamp and public key
// Prevents replay attacks across different accounts, still no state management required
// createChallenge creates a deterministic challenge from timestamp and public key
func createChallenge(timestamp int64, publicKey string) string {
	// Create a unique message combining action, timestamp, and public key
	message := fmt.Sprintf("create_account:%d:%s", timestamp, publicKey)

	// Hash the message
	hash := sha256.Sum256([]byte(message))
	return hex.EncodeToString(hash[:])
}

// @Summary creates a new account/twin
// @Description creates a new account after verifying key ownership
// @Accept  json
// @Produce  json
// @Param public_key body string true "base64 encoded public key"
// @Param signature body string true "base64 encoded signature"
// @Param timestamp body uint64 true "timestamp"
// @Success 201 {object} db.Account
// @Failure 400 {object} error
// @Failure 409 {object} db.ErrRecordAlreadyExists
// @Router /accounts/ [post]
func (s *Server) createAccountHandler(c *gin.Context) {
	var req AccountCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate public key format
	if !isValidPublicKey(req.PublicKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid public key format"})
		return
	}

	// Verify timestamp is within acceptable window
	now := time.Now()
	requestTime := time.Unix(req.Timestamp, 0)
	delta := now.Sub(requestTime)

	if delta < -MaxTimestampDelta || delta > MaxTimestampDelta {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "timestamp outside acceptable window",
			"server_time": now.Unix(),
		})
		return
	}

	// Create challenge using timestamp and public key
	challenge := createChallenge(req.Timestamp, req.PublicKey)

	// Verify signature of the challenge
	valid, err := verifySignature(req.PublicKey, challenge, req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("signature verification error: %v", err)})
		return
	}

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// Now we can create new account
	account := &db.Account{
		PublicKey: req.PublicKey,
	}

	if err := s.db.CreateAccount(account); err != nil {
		if errors.Is(err, db.ErrRecordAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "account with this public key already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// getAccountHandler retrieves an account by twin ID
// @Summary Retrieve an account by twin ID
// @Description This endpoint retrieves an account by its twin ID.
// @Tags accounts
// @Accept json
// @Produce json
// @Param twin_id path uint64 true "Twin ID of the account"
// @Success 200 {object} db.Account "Account details"
// @Failure 400 {object} gin.H "Invalid twin ID"
// @Failure 404 {object} gin.H "Account not found"
// @Router /accounts/{twin_id} [get]
func (s *Server) getAccountHandler(c *gin.Context) {
	twinID, err := strconv.ParseUint(c.Param("twin_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid twin ID"})
		return
	}

	account, err := s.db.GetAccount(twinID)
	if err != nil {
		if err == db.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// verifySignature verifies an ED25519 signature
func verifySignature(publicKeyBase64, message, signatureBase64 string) (bool, error) {
	// Decode public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false, fmt.Errorf("invalid public key format: %w", err)
	}

	// Verify public key length
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false, fmt.Errorf("invalid public key size: expected %d, got %d",
			ed25519.PublicKeySize, len(publicKeyBytes))
	}

	// Decode signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %w", err)
	}

	// Verify the signature
	return ed25519.Verify(publicKeyBytes, []byte(message), signatureBytes), nil
}

// Helper function to validate public key format
func isValidPublicKey(publicKeyBase64 string) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false
	}
	return len(publicKeyBytes) == ed25519.PublicKeySize
}

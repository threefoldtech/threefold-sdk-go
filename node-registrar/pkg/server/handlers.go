package server

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/threefoldtech/tfgrid-sdk-go/node-registrar/pkg/db"
)

const (
	PubKeySize        = 32
	MaxTimestampDelta = 2 * time.Second
)

// @Summary List farms
// @Description Get a list of farms with optional filters
// @Tags farms
// @Accept json
// @Produce json
// @Param farm_name query string false "Filter by farm name"
// @Param farm_id query int false "Filter by farm ID"
// @Param twin_id query int false "Filter by twin ID"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Results per page" default(10)
// @Success 200 {object} gin.H "List of farms"
// @Failure 400 {object} gin.H "Bad request"
// @Router /farms [get]
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

// @Summary Get farm details
// @Description Get details for a specific farm
// @Tags farms
// @Accept json
// @Produce json
// @Param farm_id path int true "Farm ID"
// @Success 200 {object} gin.H "Farm details"
// @Failure 400 {object} gin.H "Invalid farm ID"
// @Failure 404 {object} gin.H "Farm not found"
// @Router /farms/{farm_id} [get]
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

// @Summary Create new farm
// @Description Create a new farm entry
// @Tags farms
// @Accept json
// @Produce json
// @Param farm body db.Farm true "Farm creation data"
// @Success 201 {object} gin.H "Farm created successfully"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 409 {object} gin.H "Farm already exists"
// @Router /farms [post]
func (s Server) createFarmHandler(c *gin.Context) {
	var farm db.Farm

	if err := c.ShouldBindJSON(&farm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse farm info: %v", err.Error())})
		return
	}

	ensureOwner(c, farm.TwinID)
	if c.IsAborted() {
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

type UpdateFarmRequest struct {
	FarmName string `json:"farm_name" binding:"required,min=1,max=40"`
}

// @Summary Update farm
// @Description Update existing farm details
// @Tags farms
// @Accept json
// @Produce json
// @Param farm_id path int true "Farm ID"
// @Param request body UpdateFarmRequest true "Farm update data"
// @Success 200 {object} gin.H "Farm updated successfully"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 404 {object} gin.H "Farm not found"
// @Router /farms/{farm_id} [patch]
func (s Server) updateFarmsHandler(c *gin.Context) {
	var req UpdateFarmRequest
	farmID := c.Param("farm_id")

	id, err := strconv.ParseUint(farmID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid farm_id: %v", err.Error())})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse farm info: %v", err.Error())})
		return
	}

	existingFarm, err := s.db.GetFarm(id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Farm not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	ensureOwner(c, existingFarm.TwinID)
	if c.IsAborted() {
		return
	}

	// No need to hit DB if new farm name is same as the old one
	if existingFarm.FarmName != req.FarmName {
		err = s.db.UpdateFarm(id, req.FarmName)
		if err != nil {
			status := http.StatusBadRequest

			if errors.Is(err, db.ErrRecordNotFound) {
				status = http.StatusNotFound
			}

			c.JSON(status, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Farm was updated successfully",
	})
}

// @Summary List nodes
// @Description Get a list of nodes with optional filters
// @Tags nodes
// @Accept json
// @Produce json
// @Param node_id query int false "Filter by node ID"
// @Param farm_id query int false "Filter by farm ID"
// @Param twin_id query int false "Filter by twin ID"
// @Param status query string false "Filter by status"
// @Param healthy query bool false "Filter by health status"
// @Param page query int false "Page number" default(1)
// @Param size query int false "Results per page" default(10)
// @Success 200 {object} gin.H "List of nodes"
// @Failure 400 {object} gin.H "Bad request"
// @Router /nodes [get]
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

// @Summary Get node details
// @Description Get details for a specific node
// @Tags nodes
// @Accept json
// @Produce json
// @Param node_id path int true "Node ID"
// @Success 200 {object} gin.H "Node details"
// @Failure 400 {object} gin.H "Invalid node ID"
// @Failure 404 {object} gin.H "Node not found"
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

type NodeRegistrationRequest struct {
	TwinID       uint64         `json:"twin_id" binding:"required,min=1"`
	FarmID       uint64         `json:"farm_id" binding:"required,min=1"`
	Resources    db.Resources   `json:"resources" binding:"required"`
	Location     db.Location    `json:"location" binding:"required"`
	Interfaces   []db.Interface `json:"interfaces" binding:"required,min=1,dive"`
	SecureBoot   bool           `json:"secure_boot"`
	Virtualized  bool           `json:"virtualized"`
	SerialNumber string         `json:"serial_number" binding:"required"`
}

// @Summary Register new node
// @Description Register a new node in the system
// @Tags nodes
// @Accept json
// @Produce json
// @Param request body NodeRegistrationRequest true "Node registration data"
// @Success 201 {object} gin.H "Node registered successfully"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 409 {object} gin.H "Node already exists"
// @Router /nodes [post]
func (s Server) registerNodeHandler(c *gin.Context) {
	var req NodeRegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ensureOwner(c, req.TwinID)
	if c.IsAborted() {
		return
	}

	node := db.Node{
		TwinID:       req.TwinID,
		FarmID:       req.FarmID,
		Resources:    req.Resources,
		Location:     req.Location,
		Interfaces:   req.Interfaces,
		SecureBoot:   req.SecureBoot,
		Virtualized:  req.Virtualized,
		SerialNumber: req.SerialNumber,
		Approved:     false, // Default to unapproved awaiting farmer approval
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

// @Summary Report node uptime
// @Description Submit uptime report for a node
// @Tags nodes
// @Accept json
// @Produce json
// @Param node_id path int true "Node ID"
// @Param request body UptimeReportRequest true "Uptime report data"
// @Success 201 {object} gin.H "Uptime reported successfully"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 404 {object} gin.H "Node not found"
// @Router /nodes/{node_id}/uptime [post]
func (s *Server) uptimeReportHandler(c *gin.Context) {
	var req UptimeReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get node
	node, err := s.db.GetNode(req.NodeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "node not found"})
		return
	}

	ensureOwner(c, node.TwinID)
	if c.IsAborted() {
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
	Timestamp int64  `json:"timestamp" binding:"required"`
	PublicKey string `json:"public_key" binding:"required"` // base64 encoded
	// the registrar expect a signature of a message with format `timestampStr:publicKeyBase64`
	// - signature format: base64(ed25519_or_sr22519_signature)
	Signature string   `json:"signature" binding:"required"`
	Relays    []string `json:"relays,omitempty"`
	RMBEncKey string   `json:"rmb_enc_key,omitempty"`
}

// @Summary Create new account
// @Description Create a new twin account with cryptographic verification
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body AccountCreationRequest true "Account creation data"
// @Success 201 {object} db.Account "Created account details"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 409 {object} gin.H "Account already exists"
// @Router /accounts [post]
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
	// Challenge is uniquely tied to both the timestamp and public key
	// Prevents replay attacks, still no state management required
	challenge := []byte(fmt.Sprintf("%d:%s", req.Timestamp, req.PublicKey))

	// Decode public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid public key format"})
	}
	// Decode signature from base64
	signatureBytes, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid signature format: %v", err)})
	}
	// Verify signature of the challenge
	err = verifySignature(publicKeyBytes, challenge, signatureBytes)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("signature verification error: %v", err)})
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

/* // verifySignature verifies an ED25519 signature
func verifySignature(publicKey, chalange, signature []byte) (bool, error) {

	// Verify the signature
	return ed25519.Verify(publicKey, chalange, signature), nil
} */

type UpdateAccountRequest struct {
	Relays    pq.StringArray `json:"relays"`
	RMBEncKey string         `json:"rmb_enc_key"`
}

// updateAccountHandler updates an account's relays and RMB encryption key
// @Summary Update account details
// @Description Updates an account's relays and RMB encryption key
// @Tags accounts
// @Accept json
// @Produce json
// @Param twin_id path uint64 true "Twin ID of the account"
// @Param account body UpdateAccountRequest true "Account details to update"
// @Success 200 {object} gin.H "Account updated successfully"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 404 {object} gin.H "Account not found"
// @Router /accounts/{twin_id} [patch]
func (s *Server) updateAccountHandler(c *gin.Context) {
	twinID, err := strconv.ParseUint(c.Param("twin_id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid twin ID"})
		return
	}

	ensureOwner(c, twinID)
	if c.IsAborted() {
		return
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = s.db.UpdateAccount(twinID, req.Relays, req.RMBEncKey)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "account not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account updated successfully"})
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

// Helper function to validate public key format
func isValidPublicKey(publicKeyBase64 string) bool {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return false
	}
	return len(publicKeyBytes) == PubKeySize
}

// Helper function to ensure the request is from the owner
func ensureOwner(c *gin.Context, twinID uint64) {
	// Retrieve twinID set by the authMiddleware
	authTwinID, exists := c.Get("twinID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	// Safe type assertion
	authID, ok := authTwinID.(uint64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication type"})
		return
	}

	// Ensure that the retrieved twinID equals to the passed twinID
	if authID != twinID || twinID == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}
}

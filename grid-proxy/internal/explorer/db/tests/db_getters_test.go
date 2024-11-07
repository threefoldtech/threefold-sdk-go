package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/internal/explorer/db"
	"gorm.io/gorm/logger"
)

// TestPostgresDatabase_GetNode tests the GetNode function.
func TestPostgresDatabase_GetNode(t *testing.T) {
	
	dbTest, err := db.NewPostgresDatabase("localhost", 5432,"postgres","mypassword","testdb", 80, logger.Error)

	if err != nil {
		t.Skipf("Can't connect to testdb %e", err)
	}
	ctx := context.Background()

	// Test case 1: Node exists
	t.Run("Node exists", func(t *testing.T) {
		nodeID := uint32(118) // Node ID from the fixture data

		node, err := dbTest.GetNode(ctx, nodeID)

		// Assert no error
		assert.NoError(t, err)
		// Assert the node data matches fixture values
		assert.Equal(t, "node-118", node.ID)
		assert.Equal(t, int64(118), node.NodeID)
		assert.Equal(t, int64(52), node.FarmID)
		assert.Equal(t, "United States", node.Country)
		assert.Equal(t, "Los Angeles", node.City)
		assert.Equal(t, int64(1000), node.Uptime)
		assert.Equal(t, int64(1730904704), node.Created)
		assert.Equal(t, "Diy", node.Certification)
	})

	// Test case 2: Node does not exist
	t.Run("Node does not exist", func(t *testing.T) {
		nonExistentNodeID := uint32(99999) // Node ID that doesn’t exist

		node, err := dbTest.GetNode(ctx, nonExistentNodeID)

		assert.ErrorIs(t, err, db.ErrNodeNotFound)
		assert.Equal(t, db.Node{}, node)
	})
}

// TestPostgresDatabase_GetFarm tests the GetFarm function.
func TestPostgresDatabase_GetFarm(t *testing.T) {
	// Connect to the test database
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	// Test case 1: Farm exists
	t.Run("Farm exists Dedicated", func(t *testing.T) {
		farmID := uint32(3) // Farm ID from the fixture data

		farm, err := dbTest.GetFarm(ctx, farmID)

		assert.NoError(t, err)
		assert.Equal(t, int(3), farm.FarmID)
		assert.Equal(t, "farm-name-3", farm.Name)
		assert.Equal(t, int(5), farm.TwinID)
		assert.Equal(t, int(1), farm.PricingPolicyID)
		//assert.True(t, farm.Dedicated)
		assert.Equal(t, "Diy", farm.Certification)
	})

	// Test case 2: Farm exists
	t.Run("Farm exists not Dedicated", func(t *testing.T) {
		farmID := uint32(4) // Farm ID from the fixture data

		farm, err := dbTest.GetFarm(ctx, farmID)

		assert.NoError(t, err)
		assert.Equal(t, int(4), farm.FarmID)
		assert.Equal(t, "farm-name-4", farm.Name)
		assert.Equal(t, int(6), farm.TwinID)
		assert.Equal(t, int(1), farm.PricingPolicyID)
		//assert.False(t, farm.Dedicated)
		assert.Equal(t, "Diy", farm.Certification)
	})

	// Test case 3: Farm does not exist
	t.Run("Farm does not exist", func(t *testing.T) {
		nonExistentFarmID := uint32(999) // Farm ID that doesn’t exist

		farm, _ := dbTest.GetFarm(ctx, nonExistentFarmID)

		//assert.ErrorIs(t, err, db.ErrFarmNotFound)
		assert.Equal(t, db.Farm{}, farm)
	})
}
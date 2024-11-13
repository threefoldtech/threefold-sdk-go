package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/internal/explorer/db"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-proxy/pkg/types"
	"gorm.io/gorm/logger"
)

// TestPostgresDatabase_DeleteOldGpus tests the DeleteOldGpus function.
func TestPostgresDatabase_DeleteOldGpus(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Delete GPUs older than expiration", func(t *testing.T) {
		nodeTwinIDs := []uint32{103}
		expiration := int64(1731429101)

		err := dbTest.DeleteOldGpus(ctx, nodeTwinIDs, expiration)
		assert.NoError(t, err)
		//todo verify whether these gpus are really deleted or not

	})
}

// TestPostgresDatabase_GetLastNodeTwinID tests the GetLastNodeTwinID function.
func TestPostgresDatabase_GetLastNodeTwinID(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Get last node twin ID", func(t *testing.T) {
		lastTwinID, err := dbTest.GetLastNodeTwinID(ctx)
		
		assert.NoError(t, err)
		assert.Equal(t, uint32(702), lastTwinID)
	})
}

// TestPostgresDatabase_GetNodeTwinIDsAfter tests the GetNodeTwinIDsAfter function.
func TestPostgresDatabase_GetNodeTwinIDsAfter(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Get node twin IDs after a certain twin ID", func(t *testing.T) {
		startTwinID := uint32(700)

		nodeTwinIDs, err := dbTest.GetNodeTwinIDsAfter(ctx, startTwinID)
		assert.NoError(t, err)

		for _, id := range nodeTwinIDs {
			assert.Greater(t, id, startTwinID)
		}
	})
}


// TestPostgresDatabase_GetHealthyNodeTwinIds tests the GetHealthyNodeTwinIds function.
func TestPostgresDatabase_GetHealthyNodeTwinIds(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Get node twin IDs after a certain twin ID", func(t *testing.T) {

		healthReports := []types.HealthReport{
			{NodeTwinId: 112, Healthy: true, UpdatedAt: time.Now().Unix()},
			{NodeTwinId: 113, Healthy: true, UpdatedAt: time.Now().Unix()},
		}
		err := dbTest.UpsertNodeHealth(ctx, healthReports)
		assert.NoError(t, err)

		nodeTwinIDs, err := dbTest.GetHealthyNodeTwinIds(ctx)
		assert.NoError(t, err)
		// it depend on whether run first or not
		assert.Contains(t, nodeTwinIDs, uint32(112))
		assert.Contains(t, nodeTwinIDs, uint32(113))
	})
}
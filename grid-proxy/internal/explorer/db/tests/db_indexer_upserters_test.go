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

// TestPostgresDatabase_UpsertNodesGPU tests the UpsertNodesGPU function.
func TestPostgresDatabase_UpsertNodesGPU(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Nodes GPU", func(t *testing.T) {
		//take care not to interfere with other tests like (TestPostgresDatabase_DeleteOldGpus)
		gpus := []types.NodeGPU{
			{ID: "node-gpu-104-3", NodeTwinID: 104, Vendor: "NVIDIA", Device: "RTX 3090", Contract: 1, UpdatedAt: time.Now().Unix()},
			{ID: "node-gpu-102-0", NodeTwinID: 102, Vendor: "AMD", Device: "RX 6800", Contract: 1, UpdatedAt: time.Now().Unix()},
		}
		err := dbTest.UpsertNodesGPU(ctx, gpus)
		assert.NoError(t, err)
		// TODO check number of gpus for the nodes after finnish (currently node gpu doesn't map right)
		// res, err := dbTest.GetNode(ctx, uint32(104))
		// assert.NoError(t, err)
		// assert.Equal(t, len(res.Gpus), 3)

		// res, err = dbTest.GetNode(ctx, uint32(102))
		// assert.NoError(t, err)
		// assert.Equal(t, len(res.Gpus), 1)

	})
}

// TestPostgresDatabase_UpsertNodeHealth tests the UpsertNodeHealth function.
func TestPostgresDatabase_UpsertNodeHealth(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Node Health", func(t *testing.T) {
		healthReports := []types.HealthReport{
			//two ids that aren't in the health table
			{NodeTwinId: 115, Healthy: true, UpdatedAt: time.Now().Unix()},
			{NodeTwinId: 114, Healthy: false, UpdatedAt: time.Now().Unix()},
		}

		countOfHealthyNodeIds, err := dbTest.GetHealthyNodeTwinIds(ctx)
		assert.NoError(t, err)
		err = dbTest.UpsertNodeHealth(ctx, healthReports)
		assert.NoError(t, err)

		currCountOfHealthyNodeIds, err := dbTest.GetHealthyNodeTwinIds(ctx)
		assert.NoError(t, err)

		assert.Equal(t, len(countOfHealthyNodeIds)+1, len(currCountOfHealthyNodeIds))
	})
}

// TestPostgresDatabase_UpsertNodeDmi tests the UpsertNodeDmi function.
func TestPostgresDatabase_UpsertNodeDmi(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Node DMI", func(t *testing.T) {
		dmis := []types.Dmi{
			{
				NodeTwinId: 105,
				BIOS: types.BIOS{
					Vendor:  "American Megatrends",
					Version: "v1.0",
				},
				Baseboard: types.Baseboard{
					Manufacturer: "ASUS",
					ProductName:  "Prime Z390-A",
				},
				Processor: []types.Processor{
					{
						Version:     "Intel Core i7-9700K",
						ThreadCount: "8",
					},
				},
				Memory: []types.Memory{
					{
						Manufacturer: "Kingston",
						Type:         "DDR4 16GB",
					},
				},
				UpdatedAt: time.Now().Unix(),
			},
			{
				NodeTwinId: 106,
				BIOS: types.BIOS{
					Vendor:  "Phoenix Technologies",
					Version: "v2.0",
				},
				Baseboard: types.Baseboard{
					Manufacturer: "Gigabyte",
					ProductName:  "B450 AORUS PRO WIFI",
				},
				Processor: []types.Processor{
					{
						Version:     "AMD Ryzen 7 3700X",
						ThreadCount: "16",
					},
				},
				Memory: []types.Memory{
					{
						Manufacturer: "Corsair",
						Type:         "DDR4 32GB",
					},
				},
				UpdatedAt: time.Now().Unix(),
			},
		}

		err := dbTest.UpsertNodeDmi(ctx, dmis)
		assert.NoError(t, err)
		//todo verify whether these Dmi's are really upserted right or not

	})
}

// TestPostgresDatabase_UpsertNetworkSpeed tests the UpsertNetworkSpeed function.
func TestPostgresDatabase_UpsertNetworkSpeed(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Network Speed", func(t *testing.T) {
		speeds := []types.Speed{
			{NodeTwinId: 104, Download: 100.5, Upload: 50.2, UpdatedAt: time.Now().Unix()},
			{NodeTwinId: 105, Download: 150.8, Upload: 75.1, UpdatedAt: time.Now().Unix()},
		}

		err := dbTest.UpsertNetworkSpeed(ctx, speeds)
		assert.NoError(t, err)
		//todo verify whether these speed's are really upserted right or not
	})
}

// TestPostgresDatabase_UpsertNodeIpv6Report tests the UpsertNodeIpv6Report function.
func TestPostgresDatabase_UpsertNodeIpv6Report(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Node IPv6 Report", func(t *testing.T) {
		ips := []types.HasIpv6{
			{NodeTwinId: 104, HasIpv6: true, UpdatedAt: time.Now().Unix()},
			{NodeTwinId: 105, HasIpv6: false, UpdatedAt: time.Now().Unix()},
		}

		err := dbTest.UpsertNodeIpv6Report(ctx, ips)
		assert.NoError(t, err)
		//todo verify whether these HasIpv6's are really upserted right or not
	})
}

// TestPostgresDatabase_UpsertNodeWorkloads tests the UpsertNodeWorkloads function.
func TestPostgresDatabase_UpsertNodeWorkloads(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Node Workloads", func(t *testing.T) {
		workloads := []types.NodesWorkloads{
			{NodeTwinId: 101, WorkloadsNumber: 5, UpdatedAt: time.Now().Unix()},
			{NodeTwinId: 102, WorkloadsNumber: 3, UpdatedAt: time.Now().Unix()},
		}

		err := dbTest.UpsertNodeWorkloads(ctx, workloads)
		assert.NoError(t, err)
		//todo verify whether these NodesWorkloads's are really upserted right or not
	})
}

// TestPostgresDatabase_GetLastUpsertsTimestamp tests the GetLastUpsertsTimestamp function.
func TestPostgresDatabase_GetLastUpsertsTimestamp(t *testing.T) {
	dbTest, err := db.NewPostgresDatabase("localhost", 5432, "postgres", "mypassword", "testdb", 80, logger.Error)
	if err != nil {
		t.Skipf("Can't connect to testdb: %v", err)
	}
	ctx := context.Background()

	t.Run("Upsert Node Workloads", func(t *testing.T) {
		currTime := time.Now().Unix()
		workloads := []types.NodesWorkloads{
			{NodeTwinId: 104, WorkloadsNumber: 5, UpdatedAt: currTime},
		}
		err := dbTest.UpsertNodeWorkloads(ctx, workloads)
		assert.NoError(t, err)

		state, err := dbTest.GetLastUpsertsTimestamp()
		assert.NoError(t, err)
		assert.Equal(t, currTime, state.Workloads.UpdatedAt)
	})
}

package db

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"database/sql"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/logger"
)



func TestMain(m *testing.M){
	
	testDb, err := sql.Open("postgres","host=localhost user=postgres port=5432 password=mypassword dbname=testdb sslmode=disable") 

	if err != nil{
		log.Fatalf("could not fill db fixtures: %v", err)
	}
	defer testDb.Close()

	// Load and execute schema
	schema, err := os.ReadFile("../../../tools/db/schema.sql")
	if err != nil {
		log.Fatalf("could not load schema sql file: %v", err)
	}
	_, err = testDb.Exec(string(schema))
	if err != nil {
		log.Fatalf("could not apply schema: %v", err)
	}

	//
	// Load and execute schema
	setup, err := os.ReadFile("../../../tools/db/setup.sql")
	if err != nil {
		log.Fatalf("could not load setup sql file: %v", err)
	}
	_, err = testDb.Exec(string(setup))
	if err != nil {
		log.Fatalf("could not apply setup: %v", err)
	}

	// it looks like a useless block but everything breaks when it's removed
	_, err = testDb.Query("SELECT current_database();")
	if err != nil {
		panic(err)
	}

	// Load and execute fixture data
	queries, err := os.ReadFile("./fixtures/testdata.sql")
	if err != nil {
		log.Fatalf("could not load query sql file: %v", err)
	}
	_, err = testDb.Exec(string(queries))
	if err != nil {
		log.Fatalf("could not populate db: %v", err)
	}
	time.Sleep(5* time.Second)
	code := m.Run()
	// Drop the database after tests
	defer testDb.Exec(
		`
		DROP TABLE IF EXISTS account CASCADE;
		DROP TABLE IF EXISTS burn_transaction CASCADE;
		DROP TABLE IF EXISTS city CASCADE;
		DROP TABLE IF EXISTS contract_bill_report CASCADE;
		DROP TABLE IF EXISTS contract_resources CASCADE;
		DROP TABLE IF EXISTS country CASCADE;
		DROP TABLE IF EXISTS entity CASCADE;
		DROP TABLE IF EXISTS entity_proof CASCADE;
		DROP TABLE IF EXISTS farm CASCADE;
		DROP TABLE IF EXISTS farming_policy CASCADE;
		DROP TABLE IF EXISTS historical_balance CASCADE;
		DROP TABLE IF EXISTS interfaces CASCADE;
		DROP TABLE IF EXISTS location CASCADE;
		DROP TABLE IF EXISTS migrations CASCADE;
		DROP TABLE IF EXISTS mint_transaction CASCADE;
		DROP TABLE IF EXISTS name_contract CASCADE;
		DROP TABLE IF EXISTS node CASCADE;
		DROP TABLE IF EXISTS node_contract CASCADE;
		DROP TABLE IF EXISTS node_resources_free CASCADE;
		DROP TABLE IF EXISTS node_resources_total CASCADE;
		DROP TABLE IF EXISTS node_resources_used CASCADE;
		DROP TABLE IF EXISTS nru_consumption CASCADE;
		DROP TABLE IF EXISTS pricing_policy CASCADE;
		DROP TABLE IF EXISTS public_config CASCADE;
		DROP TABLE IF EXISTS public_ip CASCADE;
		DROP TABLE IF EXISTS refund_transaction CASCADE;
		DROP TABLE IF EXISTS rent_contract CASCADE;
		DROP TABLE IF EXISTS transfer CASCADE;
		DROP TABLE IF EXISTS twin CASCADE;
		DROP TABLE IF EXISTS typeorm_metadata CASCADE;
		DROP TABLE IF EXISTS uptime_event CASCADE;
		DROP SCHEMA IF EXISTS substrate_threefold_status CASCADE;
		DROP TABLE IF EXISTS node_gpu CASCADE;
		
	`)

	os.Exit(code)
	
}

// TestPostgresDatabase_GetNode tests the GetNode function.
func TestPostgresDatabase_GetNode(t *testing.T) {
	
	db, err := NewPostgresDatabase("localhost", 5432,"postgres","mypassword","testdb", 80, logger.Error)

	if err != nil {
		t.Skipf("Can't connect to testdb %e", err)
	}
	ctx := context.Background()

	// Test case 1: Node exists
	t.Run("Node exists", func(t *testing.T) {
		nodeID := uint32(118) // Node ID from the fixture data

		node, err := db.GetNode(ctx, nodeID)

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

		node, err := db.GetNode(ctx, nonExistentNodeID)

		// Assert error is ErrNodeNotFound
		assert.ErrorIs(t, err, ErrNodeNotFound)
		// Assert returned node is empty
		assert.Equal(t, Node{}, node)
	})
}
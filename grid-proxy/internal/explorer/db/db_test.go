package db

import (
	"log"
	"testing"

	"database/sql"

	"github.com/tanimutomo/sqlfile"

)


var db PostgresDatabase

func TestMain(m *testing.T){
	
	testDb, err := sql.Open("postgres","host=localhost user=postgres port=5432 password=mypassword dbname=testDB sslmode=disable") 

	if err != nil{
		log.Fatalf("could not fill db fixtures: %v", err)
	}

	s := sqlfile.New()

	err = s.Files("testdata.sql")

	if err!= nil{
		log.Fatalf("could not load sql file: %v",err)
	}
	_, err = s.Exec(testDb)

	if err !=nil{
		log.Fatalf("could not fill the db: %v", err)
	}
}

package data

import (
	"log"

	"github.com/hashicorp/go-memdb"
)

// SetupDB uses the schema to validate and create MemDB
func SetupDB(schema *memdb.DBSchema) *memdb.MemDB {
	err := schema.Validate()
	if err != nil {
		log.Fatal(err)
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

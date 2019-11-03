package write

import (
	"log"

	"github.com/hashicorp/go-memdb"
)

// Schema is the schema for the write model
var Schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"product": &memdb.TableSchema{
			Name: "product",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "SKU"},
				},
			},
		},
	},
}

// SeedDB seeds the db with data
func SeedDB(db *memdb.MemDB) {
	log.Println("seeding products...")
	txn := db.Txn(true)
	defer txn.Abort()

	err := txn.Insert("product", &product{"Test1", "ABC232323", 100})
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Insert("product", &product{"Test2", "ABC883388", 100})
	if err != nil {
		log.Fatal(err)
	}

	txn.Commit()
	log.Println("seeded products")
}

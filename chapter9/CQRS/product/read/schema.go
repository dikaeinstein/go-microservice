package read

import (
	"github.com/hashicorp/go-memdb"
)

// Schema is the schema for the read model
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

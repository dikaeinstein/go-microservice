package read

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hashicorp/go-memdb"
)

// Handler represents the read http handler
type Handler struct {
	db *memdb.MemDB
}

// NewHandler creates a new read Handler
func NewHandler(db *memdb.MemDB) *Handler {
	return &Handler{db}
}

func (h *Handler) GetProducts(rw http.ResponseWriter, r *http.Request) {
	log.Println("/get handler called")

	txn := h.db.Txn(false)
	defer txn.Abort()

	results, err := txn.Get("product", "id")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	products := make([]product, 0)
	for {
		obj := results.Next()
		if obj == nil {
			break
		}

		products = append(products, obj.(product))
	}

	rw.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.Encode(products)
}

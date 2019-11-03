package write

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-memdb"
	"github.com/nats-io/nats.go"
)

// Handler represents the write http handler
type Handler struct {
	natsClient *nats.Conn
	db         *memdb.MemDB
}

// NewHandler creates a new write Handler
func NewHandler(db *memdb.MemDB, natsClient *nats.Conn) *Handler {
	return &Handler{natsClient, db}
}

// InsertProduct writes a new product to event store and publishes the event
func (h *Handler) InsertProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("/insert handler called")

	p := &product{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	txn := h.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("product", p); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	txn.Commit()

	w.WriteHeader(http.StatusNoContent)
	h.natsClient.Publish("product.inserted", data)
}

// StockCount returns the stock count for a product
func (h *Handler) StockCount(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	txn := h.db.Txn(false)
	defer txn.Abort()

	obj, err := txn.First("product", "id", id)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if obj == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	p := obj.(*product)
	rw.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(rw, `{"quantity": %v}`, p.StockCount)
}

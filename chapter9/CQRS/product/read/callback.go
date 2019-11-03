package read

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/go-memdb"
	"github.com/nats-io/nats.go"
)

// MakeProductMessageCallBack is a factory function that creates the product message callback
// that is used as the asynchronous subscriber by NATS
func MakeProductMessageCallBack(db *memdb.MemDB) nats.MsgHandler {
	return func(m *nats.Msg) {
		log.Println("new event")
		pie := productInsertedEvent{}
		err := json.Unmarshal(m.Data, &pie)
		if err != nil {
			log.Println("Unable to unmarshal event object")
			return
		}

		p := product{}.FromProductInsertedEvent(pie)

		txn := db.Txn(true)
		if err := txn.Insert("product", p); err != nil {
			log.Println(err)
			return
		}
		txn.Commit()

		log.Println("Saved product: ", p)
	}
}

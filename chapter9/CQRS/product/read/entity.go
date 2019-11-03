package read

import "time"

type productInsertedEvent struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

type product struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	LastUpdated string `json:"last_updated"`
}

func (p product) FromProductInsertedEvent(e productInsertedEvent) product {
	p.Name = e.Name
	p.SKU = e.SKU
	p.LastUpdated = time.Now().Format(time.RFC3339)

	return p
}

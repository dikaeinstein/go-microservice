package write

type product struct {
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	StockCount int    `json:"stock_count"`
}

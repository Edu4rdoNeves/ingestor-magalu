package dto

type PulseData struct {
	Tenant     string  `json:"tenant"`
	ProductSku string  `json:"product_sku"`
	UseUnity   string  `json:"use_unity"`
	UsedAmount float64 `json:"used_amount"`
}

package dto

type PulseData struct {
	Tenant     string  `json:"tenant"`
	ProductSku string  `json:"product_sku"`
	UseUnity   string  `json:"use_unity"`
	UsedAmount float64 `json:"used_amount"`
}

type PopulateQueueParams struct {
	TotalMessages int `json:"total_messages"`
	WorkersNumber int `json:"workers_number"`
	BufferSize    int `json:"buffer_size"`
}

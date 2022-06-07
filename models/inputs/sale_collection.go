package inputs

type SaleCollection struct {
	SaleID          string `json:"sale_id"`
	NodeURL         string `json:"node_url"`
	TransactionHash string `json:"transaction_hash"`
}

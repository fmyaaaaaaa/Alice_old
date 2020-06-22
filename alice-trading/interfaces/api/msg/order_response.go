package msg

type OrderResponse struct {
	OrderCreateTransaction Transaction          `json:"orderCreateTransaction"`
	OrderFillTransaction   OrderFillTransaction `json:"orderFillTransaction"`
	RelatedTransactionIDs  []string             `json:"relatedTransactionIDs"`
	LastTransactionID      string               `json:"lastTransactionID"`
}

package msg

type OrderErrorResponse struct {
	OrderRejectTransaction Transaction `json:"orderRejectTransaction"`
	RelatedTransactionIDs  []string    `json:"relatedTransactionIDs"`
	LastTransactionID      string      `json:"lastTransactionID"`
	ErrorCode              string      `json:"errorCode"`
	ErrorMessage           string      `json:"errorMessage"`
}

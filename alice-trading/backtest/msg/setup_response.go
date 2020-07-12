package msg

type SetupResponse struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

func NewSetupResponse(status int, result string) SetupResponse {
	return SetupResponse{
		Status: status,
		Result: result,
	}
}

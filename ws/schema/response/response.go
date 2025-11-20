package response

type ResponseBody struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"requestId,omitempty"`
}

type BroadcastBody struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

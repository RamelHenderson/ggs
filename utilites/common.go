package utilites

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    any         `json:"meta,omitempty"`
}

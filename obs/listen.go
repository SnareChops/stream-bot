package obs

type OBSResponse struct {
	ID     string `json:"message-id"`
	Status string `json:"status"` // Response status, will be "ok" or "error"
	Error  string `json:"error,omitempty"`
}

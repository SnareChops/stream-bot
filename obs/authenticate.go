package obs

type OBSGetAuthResponse struct {
	ID           string `json:"message-id"`
	Status       string `json:"status"` // Response status, will be "ok" or "error"
	Error        string `json:"error,omitempty"`
	AuthRequired bool   `json:"authRequired"`
	Challenge    string `json:"challenge"`
	Salt         string `json:"salt"`
}

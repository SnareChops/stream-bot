package signals

type OBSRequest struct {
	Type string
	Args map[string]any
}

func NewOBSRequest(typ string, args map[string]any) OBSRequest {
	return OBSRequest{typ, args}
}

var (
	Shutdown         chan bool
	SendToUI         chan []byte
	SendToAdmin      chan []byte
	SendToIRC        chan string
	SendToOBS        chan OBSRequest
	ReceiveFromAdmin chan []byte
)

func init() {
	Shutdown = make(chan bool)
	SendToUI = make(chan []byte, 20)
	SendToAdmin = make(chan []byte, 20)
	SendToIRC = make(chan string)
	SendToOBS = make(chan OBSRequest)
	ReceiveFromAdmin = make(chan []byte, 20)
}

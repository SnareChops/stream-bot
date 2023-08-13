package signals

var (
	Shutdown    chan bool
	SendToUI    chan []byte
	SendToAdmin chan []byte
	SendToIRC   chan string
)

func init() {
	Shutdown = make(chan bool)
	SendToUI = make(chan []byte, 20)
	SendToAdmin = make(chan []byte, 20)
	SendToIRC = make(chan string)
}

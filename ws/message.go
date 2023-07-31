package ws

const (
	MessageTypeWelcome      = "session_welcome"
	MessageTypeKeepalive    = "session_keepalive"
	MessageTypeNotification = "notification"
	MessageTypeReconnect    = "session_reconnect"
	MessageTypeRevocation   = "revocation"
)

type Message struct {
	Metadata MessageMetadata        `json:"metadata"`
	Payload  map[string]interface{} `json:"payload"`
}

type MessageMetadata struct {
	MessageID           string `json:"message_id"`
	MessageType         string `json:"message_type"`
	MessageTimestamp    string `json:"message_timestamp"`
	SubscriptionType    string `json:"subscription_type"`
	SubscriptionVersion string `json:"subscription_version"`
}

type SessionPayload struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	ConnectedAt      string `json:"connected_at"`
	KeepaliveTimeout int    `json:"keepalive_timeout_seconds"`
	ReconnectUrl     string `json:"reconnect_url"`
}

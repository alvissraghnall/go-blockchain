package node

// MessageType represents the type of a network message.
type MessageType string

const (
	MessageTypeTransaction MessageType = "TRANSACTION"
	MessageTypeBlock       MessageType = "BLOCK"
)

// Message represents a network message.
type Message struct {
	Type    MessageType
	Payload string
}

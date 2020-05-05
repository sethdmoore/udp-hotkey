package types

// Packet is sent over the Serial interface
type Packet struct {
	Action    uint8
	Modifiers uint16
	KeyCode   uint8
}

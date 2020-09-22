package types

// Packet is a struct encoded/decoded with gob
type Packet struct {
	Action    uint8
	Modifiers uint16
	KeyCode   uint8
}

package constants

const (
	// Press for single hotkey press
	KeyPress = uint8(iota)
	// Held for hotkey held
	KeyHeld
	// Release corresponds to Held, Release hotkey
	KeyRelease

	// Byte size for packet
	PacketLength = 4
)

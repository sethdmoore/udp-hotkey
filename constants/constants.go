package constants

const (
	// Press for single hotkey press
	KeyPress = uint8(iota)
	// Held for hotkey held
	KeyHeld
	// Release corresponds to Held, Release hotkey
	KeyRelease
)

const (
	ModAlt = 1 << uint16(iota)
	ModCtrl
	ModShift
	ModWin
	ModNoRepeat = uint16(16384)
)

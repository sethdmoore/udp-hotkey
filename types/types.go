package types

// Packet is sent over the Serial interface
type Packet struct {
	Action    uint8
	Modifiers uint16
	KeyCode   uint8
}

// God bless Vladimir Vivien. This article was paramount in my understanding
// of serializing bytes over the wire. Not only does the author show you the
// manual process, they also show you the easy way with structs afterwards.

// https://medium.com/learning-the-go-programming-language/encoding-data-with-the-go-binary-package-42c7c0eb3e73
// https://archive.vn/jY7tm

package domain

// EthernetConnection behinaltet die Verbindungen
type EthernetConnection struct {
	Src             string
	Dst             string
	EthernetType    string
	NumberOfPackets int
}

// Connections Map for each src dst and type combination with the number of connectins
var Conncetions = make(map[string]EthernetConnection)

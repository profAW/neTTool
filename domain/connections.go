package domain

// EthernetConnection behinaltet die Verbindungen
type EthernetConnection struct {
	Src             string
	Dst             string
	EthernetType    string
	NumberOfPackets int
}

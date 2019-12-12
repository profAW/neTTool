package usecases

import (
	"github.com/google/gopacket"
)

// PackagePort gets the RawData Packges
type PackagePort interface {
	Read() map[int]gopacket.Packet
}

// UcGetNetworkData Usecase to read network data from source
type UcGetNetworkData struct {
	Source PackagePort
}

func (e UcGetNetworkData) Read() map[int]gopacket.Packet {
	return e.Source.Read()
}

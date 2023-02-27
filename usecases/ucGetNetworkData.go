package usecases

import (
	"github.com/google/gopacket"
	"neTTool/domain"
)

// PackagePort gets the RawData Packges
type PackagePort interface {
	Read() map[int]gopacket.Packet
}

// UcGetNetworkData Use-Case to read network data from source
type UcGetNetworkData struct {
	Source PackagePort
}

func (e UcGetNetworkData) Read() map[int]gopacket.Packet {
	return e.Source.Read()
}

func (e UcGetNetworkData) CreateNetworkData(Data map[int]gopacket.Packet) (map[string]domain.CommonConnection, map[string]domain.Node) {
	return domain.CreateConnectionList(Data)
}

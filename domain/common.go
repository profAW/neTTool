package domain

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strconv"
	"time"
)

type CommonConnection struct {
	Src             string
	Dst             string
	EthernetType    string
	NumberOfPackets int
	Ts              []time.Time
	DeltaTS         []float64
}

func getEthernetTyp(etype []byte) string {
	// Extract EthernetType from Byte-String on Pos 12 and 13
	v1 := int64(etype[12])
	v2 := int64(etype[13])
	return strconv.FormatInt(v1, 16) + strconv.FormatInt(v2, 16)
}

func GetKeyAndEthernetTyp(ethernetPacket *layers.Ethernet) (string, string) {

	flow := ethernetPacket.LinkFlow()
	etherTyp := ""

	if ethernetPacket.EthernetType.String() == "UnknownEthernetType" {
		etherTyp = getEthernetTyp(ethernetPacket.LayerContents())

	} else {
		etherTyp = ethernetPacket.EthernetType.String()
	}

	key := flow.Dst().String() + "->" + flow.Src().String() + "|" + etherTyp

	return key, etherTyp
}

func CreateConnectionList(Data map[int]gopacket.Packet) map[string]CommonConnection {

	var connection = make(map[string]CommonConnection) // Map mit flow connections and number of connections

	for _, packet := range Data {

		//fmt.Println(packet)
		// Let's see if the packet is an ethernet packet
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			//fmt.Println("Ethernet layer detected.")
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

			var key, etherTyp = GetKeyAndEthernetTyp(ethernetPacket)

			_, ok := connection[key]
			if ok {
				con := connection[key]
				con.NumberOfPackets++
				if etherTyp == "8892" {
					con.Ts = append(con.Ts, packet.Metadata().Timestamp)
				}
				connection[key] = con

			} else {

				var con CommonConnection
				con.EthernetType = etherTyp
				con.Src = ethernetPacket.SrcMAC.String()
				con.Dst = ethernetPacket.DstMAC.String()
				con.Ts = append(con.Ts, packet.Metadata().Timestamp)
				con.NumberOfPackets = 1
				connection[key] = con
			}
		}

	}
	return connection
}

package domain

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strconv"
	"time"
)

const pnEtherType = "8892"

type CommonConnection struct {
	MacSrc          string
	MacDst          string
	EthernetType    string
	NumberOfPackets int
	Ts              []time.Time
	DeltaTS         []float64
	IPSrc           string
	IPDst           string
	PortSrc         string
	PortDst         string
	SocketType      string
}

func (e CommonConnection) GetKey() string {
	return e.MacDst + "->" + e.MacSrc + "|" + e.EthernetType
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

func GetLayer4Key(packet gopacket.Packet) string {

	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	ipPacket := ipLayer.(*layers.IPv4)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	tcpLayer := packet.Layer(layers.LayerTypeTCP)

	key := ipPacket.SrcIP.String() + "->" + ipPacket.DstIP.String()

	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		key += " | UDP " + udp.SrcPort.String() + " --> " + udp.DstPort.String()
	}
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		key += " | UDP " + tcp.SrcPort.String() + " --> " + tcp.DstPort.String()
	}

	return key
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

			var key, etherType = GetKeyAndEthernetTyp(ethernetPacket)

			_, ok := connection[key]
			if ok {
				con := connection[key]
				con.NumberOfPackets++

				if etherType == pnEtherType {
					con.Ts = append(con.Ts, packet.Metadata().Timestamp)
				}
				connection[key] = con

			} else {

				var con CommonConnection
				con.EthernetType = etherType
				con.MacSrc = ethernetPacket.SrcMAC.String()
				con.MacDst = ethernetPacket.DstMAC.String()
				con.Ts = append(con.Ts, packet.Metadata().Timestamp)
				con.NumberOfPackets = 1
				connection[key] = con
			}

			ipLayer := packet.Layer(layers.LayerTypeIPv4)

			if ipLayer != nil {
				var key2 = GetLayer4Key(packet)

				_, ok2 := connection[key2]

				if ok2 {
					con := connection[key2]
					con.NumberOfPackets++
					connection[key2] = con

				} else {
					var con CommonConnection

					con.MacSrc = ethernetPacket.SrcMAC.String()
					con.MacDst = ethernetPacket.DstMAC.String()
					con.NumberOfPackets = 1

					ipLayer := packet.Layer(layers.LayerTypeIPv4)
					ipPacket := ipLayer.(*layers.IPv4)

					con.IPSrc = ipPacket.SrcIP.String()
					con.IPDst = ipPacket.DstIP.String()

					udpLayer := packet.Layer(layers.LayerTypeUDP)
					tcpLayer := packet.Layer(layers.LayerTypeTCP)

					if udpLayer != nil {
						con.SocketType = "UDP"
						udp, _ := udpLayer.(*layers.UDP)
						con.PortDst = udp.DstPort.String()
						con.PortSrc = udp.SrcPort.String()
					}
					if tcpLayer != nil {
						con.SocketType = "TCP"
						tcp, _ := tcpLayer.(*layers.TCP)
						con.PortDst = tcp.DstPort.String()
						con.PortSrc = tcp.SrcPort.String()
					}
					con.EthernetType = etherType + "|" + con.SocketType
					connection[key2] = con
				}
			}

		}

	}
	return connection
}

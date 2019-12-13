package usecases

import (
	"github.com/emicklei/dot"
	"github.com/google/gopacket"

	"strconv"

	"neTTool/domain"

	"github.com/google/gopacket/layers"
)

//ExportConnectionGraphPort export connection graph to some output
type ExportConnectionGraphPort interface {
	ExportConnectionGraph(conncetionGraph string)
}

// UcConnectionAnalysis Usecase to read network data from source
type UcConnectionAnalysis struct {
	Destination ExportConnectionGraphPort
}

var connection = make(map[string]domain.EthernetConnection) // Map mit flow connections and number of connections

// CreateConnectionList  Creation of Network-Data graph
func (e UcConnectionAnalysis) CreateConnectionList(Data map[int]gopacket.Packet) map[string]domain.EthernetConnection {

	for _, packet := range Data {

		//fmt.Println(packet)
		// Let's see if the packet is an ethernet packet
		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethernetLayer != nil {
			//fmt.Println("Ethernet layer detected.")
			ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

			flow := ethernetPacket.LinkFlow()

			etherTyp := ""

			if ethernetPacket.EthernetType.String() == "UnknownEthernetType" {

				etype := ethernetPacket.LayerContents()
				v1 := int64(etype[12])
				v2 := int64(etype[13])
				etherTyp = strconv.FormatInt(v1, 16) + strconv.FormatInt(v2, 16)
				//fmt.Println(value)
			} else {
				etherTyp = ethernetPacket.EthernetType.String()
			}

			key := flow.Dst().String() + "->" + flow.Src().String() + " | " + etherTyp //+ " - " + ethernetPacket.EthernetType.String()

			_, ok := connection[key]
			if ok {
				con := connection[key]
				con.NumberOfPackets++
				connection[key] = con
				//fmt.Println("value: ", value)
				//fmt.Println("-----")
			} else {

				var con domain.EthernetConnection
				con.EthernetType = etherTyp
				con.Src = ethernetPacket.SrcMAC.String()
				con.Dst = ethernetPacket.DstMAC.String()

				con.NumberOfPackets = 1
				connection[key] = con
			}
		}

	}
	return connection
}

// MakeConnetionGraph - create graphiv connection graph
func (e UcConnectionAnalysis) MakeConnetionGraph(connections map[string]domain.EthernetConnection) string {

	g := dot.NewGraph(dot.Directed)
	arp := g.Subgraph("ARP", dot.ClusterOption{})
	ip4 := g.Subgraph("IPv4", dot.ClusterOption{})
	ip6 := g.Subgraph("IPv6", dot.ClusterOption{})
	pn := g.Subgraph("PN", dot.ClusterOption{})
	red := g.Subgraph("Red", dot.ClusterOption{})
	notdef := g.Subgraph("Notdef", dot.ClusterOption{})

	for _, value := range connection {

		p := &arp
		switch value.EthernetType {
		case "ARP":
			p = &arp
		case "IPv4":
			p = &ip4
		case "IPv6":
			p = &ip6
		case "8892":
			p = &pn
		case "88e3":
			p = &red
		default:
			p = &notdef
		}

		res := *p
		n1 := res.Node(value.Src)
		n2 := res.Node(value.Dst)

		res.Edge(n1, n2).Attr("label", value.EthernetType)
	}
	return g.String()
}

// ExportConnectionGraph so a destination
func (e UcConnectionAnalysis) ExportConnectionGraph(conncetionGraph string) {
	e.Destination.ExportConnectionGraph(conncetionGraph)
}

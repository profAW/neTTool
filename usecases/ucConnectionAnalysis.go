package usecases

import (
	"github.com/emicklei/dot"
	"neTTool/domain"
	"strconv"
	"strings"
)

//ExportConnectionGraphPort export connection graph to some output
type ExportConnectionGraphPort interface {
	ExportConnectionGraph(conncetionGraph string)
}

// UcConnectionAnalysis Use-Case to read network data from source
type UcConnectionAnalysis struct {
	Destination ExportConnectionGraphPort
}

// MakeConnetionGraph - create graphiv connection graph
func (e UcConnectionAnalysis) MakeConnetionGraph(connections map[string]domain.CommonConnection) string {

	g := dot.NewGraph(dot.Directed)
	arp := g.Subgraph("ARP", dot.ClusterOption{})
	ip4 := g.Subgraph("IPv4", dot.ClusterOption{})
	ip4_udp := g.Subgraph("IPv4-UDP", dot.ClusterOption{})
	ip4_tcp := g.Subgraph("IPv4-TCP", dot.ClusterOption{})
	ip6 := g.Subgraph("IPv6", dot.ClusterOption{})
	pn := g.Subgraph("PN", dot.ClusterOption{})
	red := g.Subgraph("Red", dot.ClusterOption{})
	eapol := g.Subgraph("EAPOL", dot.ClusterOption{})
	linklayerdiscovery := g.Subgraph("LinkLayerDiscovery", dot.ClusterOption{})
	notdef := g.Subgraph("Notdef", dot.ClusterOption{})

	for _, value := range connections {

		p := &arp
		switch value.EthernetType {
		case "ARP":
			p = &arp
		case "IPv4":
			p = &ip4
		case "IPv4|UDP":
			p = &ip4_udp
		case "IPv4|TCP":
			p = &ip4_tcp
		case "IPv6":
			p = &ip6
		case "8892":
			p = &pn
		case "88e3":
			p = &red
		case "EAPOL":
			p = &eapol
		case "LinkLayerDiscovery":
			p = &linklayerdiscovery
		default:
			p = &notdef
		}

		res := *p

		switch value.EthernetType {
		case "IPv4|UDP":
			n1 := res.Node(value.MacSrc + " \n " + value.IPSrc)
			n2 := res.Node(value.MacDst + " \n " + value.IPDst)

			res.Edge(n1, n2).Attr("label", value.PortSrc+" \n-> \n"+value.PortDst+" \n "+strconv.Itoa(value.NumberOfPackets))

		case "IPv4|TCP":
			n1 := res.Node(value.MacSrc + " \n " + value.IPSrc)
			n2 := res.Node(value.MacDst + " \n " + value.IPDst)

			res.Edge(n1, n2).Attr("label", value.PortSrc+" \n- \n>"+value.PortDst+" \n "+strconv.Itoa(value.NumberOfPackets))
		default:
			n1 := res.Node(value.MacSrc)
			n2 := res.Node(value.MacDst)

			res.Edge(n1, n2).Attr("label", strconv.Itoa(value.NumberOfPackets))
		}
	}
	var myGraph = g.String()
	myGraph = strings.ReplaceAll(myGraph, "graph  {", "graph  {\n rankdir= LR;\n")

	return myGraph
}

// ExportConnectionGraph so a destination
func (e UcConnectionAnalysis) ExportConnectionGraph(conncetionGraph string) {
	e.Destination.ExportConnectionGraph(conncetionGraph)
}

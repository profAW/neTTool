package usecases

import (
	"github.com/emicklei/dot"
	"neTTool/domain"
	"strconv"
	"strings"
)

// ExportConnectionGraphPort export connection graph to some output
type ExportConnectionGraphPort interface {
	ExportConnectionGraph(conncetionGraph string)
}

// ExportNodeGraphPort export connection graph to some output
type ExportNodeGraphPort interface {
	ExportNodeGraph(nodes []string)
}

// UcPreparationConnections Use-Case to read network data from source
type UcPreparationConnections struct {
	Destination ExportConnectionGraphPort
}

type UcPreparationNodes struct {
	Destination ExportNodeGraphPort
}

func (e UcPreparationConnections) DoExport(connections map[string]domain.CommonConnection) {
	result := e.MakeConnetionGraph(connections)
	e.ExportConnectionGraph(result)
}

func (e UcPreparationNodes) DoExport(nodes map[string]domain.Node) {
	result := e.MakeNodeGraph(nodes)
	e.ExportNodeGraph(result)
}

// MakeConnetionGraph - create graphiv connection graph
func (e UcPreparationConnections) MakeConnetionGraph(connections map[string]domain.CommonConnection) string {

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
		case "IPv4":
			n1 := res.Node(value.MacSrc + " \n " + value.IPSrc)
			n2 := res.Node(value.MacDst + " \n " + value.IPDst)

			res.Edge(n1, n2).Attr("label", strconv.Itoa(value.NumberOfPackets))

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

func (e UcPreparationNodes) MakeNodeGraph(nodes map[string]domain.Node) []string {

	var i int
	i = 0
	var myGraph []string
	myGraph = append(myGraph, "@startuml \r")
	myGraph = append(myGraph, "  nwdiag { \r")
	myGraph = append(myGraph, "	network test_network { \r")

	var multicast []string

	for _, value := range nodes {

		if strings.HasPrefix(value.Mac, "ff:ff") || strings.HasPrefix(value.Mac, "01:00:5e") {
			multicast = append(multicast, " ' Multicast: node [address=  \" "+value.Mac+" , "+value.IP+" \"] \r")
		} else {
			myGraph = append(myGraph, "          node_"+strconv.Itoa(i)+"[address=  \" "+value.Mac+" , "+value.IP+" \"] \r")
			i++
		}
	}

	myGraph = append(myGraph, "  }\r")
	myGraph = append(myGraph, "' Multicast-Nodes\r")
	for _, value := range multicast {
		myGraph = append(myGraph, value)
	}
	myGraph = append(myGraph, "@enduml")

	return myGraph
}

// ExportConnectionGraph so a destination
func (e UcPreparationConnections) ExportConnectionGraph(conncetionGraph string) {
	e.Destination.ExportConnectionGraph(conncetionGraph)
}
func (e UcPreparationNodes) ExportNodeGraph(nodes []string) {
	e.Destination.ExportNodeGraph(nodes)
}

package usecases

import (
	"neTTool/domain"
	"sort"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// UcProfiNETAnalysis PN-Data-Analysis
type UcProfiNETAnalysis struct {
}

var conncetionPN = make(map[string]domain.ProfinetData) // Map mit flow connections and number of connections

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// GetProfiNetData  Read PN Data from Packets
func (e UcProfiNETAnalysis) GetProfiNetData(Data map[int]gopacket.Packet) map[string]domain.ProfinetData {

	i := 0
	for _, packet := range Data {

		i = i + 1

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
			}

			if etherTyp == "8892" {
				key := flow.Dst().String() + "->" + flow.Src().String() // + " | " + etherTyp //+ " - " + ethernetPacket.EthernetType.String()

				_, ok := conncetionPN[key]
				if ok {
					con := conncetionPN[key]
					//i := len(con.Ts) + 1
					//fmt.Println(i)
					con.Ts = append(con.Ts, packet.Metadata().Timestamp)
					conncetionPN[key] = con
					//fmt.Println("value: ", value)
					//fmt.Println("-----")
				} else {

					var con domain.ProfinetData
					con.Src = ethernetPacket.SrcMAC.String()
					con.Dst = ethernetPacket.DstMAC.String()
					con.Ts = append(con.Ts, packet.Metadata().Timestamp)
					conncetionPN[key] = con
				}
			}
		}
	}
	return conncetionPN
}

// CalcProfiNetDeltaTimeInMS Caluclate the TimeDiffernce between two PN packages
func (e UcProfiNETAnalysis) CalcProfiNetDeltaTimeInMS(Data map[string]domain.ProfinetData) map[string]domain.ProfinetData {

	for k, con := range Data {
		con.DeltaTS = getDeltaTimeInMs(con.Ts)
		Data[k] = con
	}
	return Data
}

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }

func getDeltaTimeInMs(timestamps []time.Time) []float64 {
	var deltaTimeInMs []float64

	var timeVektor timeSlice = timestamps

	sort.Sort(timeVektor)

	var lastTs time.Time
	first := true

	for _, value := range timeVektor {
		if !first {
			delta := value.Sub(lastTs)

			ms := float64(delta) / 1000 / 1000
			if (ms < 20) && (ms > 0.0) {

				deltaTimeInMs = append(deltaTimeInMs, ms) // / time.Millisecond)
			}
		}
		first = false
		lastTs = value

	}
	//fmt.Println(deltaTimeInMs)
	return deltaTimeInMs
}

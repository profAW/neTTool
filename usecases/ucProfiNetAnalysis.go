package usecases

import (
	"neTTool/domain"
	"sort"
	"time"
)

// UcProfiNETAnalysis PN-Data-Analysis
type UcProfiNETAnalysis struct {
}

var conncetionPN = make(map[string]domain.ProfinetConnection) // Map mit flow connections and number of connections

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CalcProfiNetDeltaTimeInMS Caluclate the TimeDiffernce between two PN packages
func (e UcProfiNETAnalysis) CalcProfiNetDeltaTimeInMS(Data map[string]domain.CommonConnection) map[string]domain.CommonConnection {

	for k, con := range Data {
		if con.EthernetType == "8892" {
			con.DeltaTS = getDeltaTimeInMs(con.Ts)
			Data[k] = con
		}
	}
	return Data
}

type timeSlice []time.Time

func (s timeSlice) Less(i, j int) bool { return s[i].Before(s[j]) }
func (s timeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s timeSlice) Len() int           { return len(s) }

func getDeltaTimeInMs(timestamps []time.Time) []float64 {
	var deltaTimeInMs []float64

	var timeVector timeSlice = timestamps

	sort.Sort(timeVector)

	var lastTs time.Time
	first := true

	for _, value := range timeVector {
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

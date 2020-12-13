package usecases

import (
	"neTTool/domain"
	"sort"
	"time"
)

// UcProfiNETAnalysis PN-Data-Analysis
type UcProfiNETAnalysis struct {
}

// CalcProfiNetDeltaTimeInMS Calculate the TimeDifference between two PN packages
func (e UcProfiNETAnalysis) CalcProfiNetDeltaTimeInMS(Data map[string]domain.CommonConnection) map[string]domain.CommonConnection {

	c := make(chan domain.CommonConnection, 100) // create non blocking channels
	numberOfOpenChannels := 0                    // count the number of open channels
	for _, con := range Data {                   // create goroutine for each connection to calculate delta ts
		go CalcDetlaTs(con, c)
		numberOfOpenChannels++
	}

	for con := range c { // loop over the results and add it to the connection map
		Data[con.GetKey()] = con
		numberOfOpenChannels--
		if numberOfOpenChannels == 0 {
			close(c)
		}
	}
	return Data
}

func CalcDetlaTs(con domain.CommonConnection, result chan domain.CommonConnection) {
	if con.EthernetType == "8892" {
		con.DeltaTS = getDeltaTimeInMs(con.Ts)
	}
	result <- con
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

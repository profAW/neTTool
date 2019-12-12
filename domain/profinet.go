package domain

import "time"

// ProfinetData contains all PN Data of one SRC-DST Connection
type ProfinetData struct {
	Src     string
	Dst     string
	Ts      []time.Time
	DeltaTS []float64
}

// Connections Map for each src dst and type combination with the number of connectins
//var Connections = make(map[string]ProfinetData)

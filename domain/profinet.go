package domain

import "time"

// ProfinetData contains all PN Data of one SRC-DST Connection
type ProfinetData struct {
	Src     string
	Dst     string
	Ts      []time.Time
	DeltaTS []float64
}

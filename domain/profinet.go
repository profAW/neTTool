package domain

import "time"

// ProfinetConnection contains all PN Data of one SRC-DST Connection
type ProfinetConnection struct {
	Src     string
	Dst     string
	Ts      []time.Time
	DeltaTS []float64
}

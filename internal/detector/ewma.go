package detector

import (
	"math"
)
type EwmaDetector struct {
	Lambda float64
	StdDev float64
	Mean float64
	Threshold float64
	Current float64
}

func (e *EwmaDetector)EwmaUpdate(value float64) (float64,bool){
	e.Current = (e.Lambda *value) + (e.Current*(1-e.Lambda))
	deviation :=math.Abs(e.Current-e.Mean)/e.StdDev
	return e.Current,deviation>e.Threshold
}
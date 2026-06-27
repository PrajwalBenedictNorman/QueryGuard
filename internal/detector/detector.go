package detector

import (
	"sync"
)

type Detector struct{
	Cusum CusumDetector
	Ewma EwmaDetector
	CusumValue float64
	CusumAlert bool
	EwmaValue float64
	EwmaAlert bool
	DriftAlert bool
}

type result struct {
	score float64
	alert bool
	methodName string
}

 func (d *Detector) CusumResult(value float64,wg *sync.WaitGroup,ch chan result ){

	score,alert := d.Cusum.CusumUpdate(value)
	ch <- result{score: score,alert: alert,methodName: "cusum"}
	wg.Done()
}	


func (d *Detector) EwmaResult(value float64 ,wg *sync.WaitGroup,ch chan result){

	score,alert := d.Ewma.EwmaUpdate(value)
	ch <- result{score: score,alert: alert,methodName: "ewma"}
	wg.Done()
}

func (d *Detector) Detection(value float64)Detector{
	ch :=make(chan result,2)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go d.CusumResult(value,wg,ch)
	go d.EwmaResult(value,wg,ch)
	wg.Wait()
	for i:=0;i<2;i++{
	r := <-ch
	if (r.methodName=="cusum"){
		d.CusumValue=r.score
		d.CusumAlert=r.alert
	}else{
		d.EwmaValue=r.score
		d.EwmaAlert=r.alert
	}
	}
	if d.CusumAlert || d.EwmaAlert {
		d.DriftAlert=true
	}
	return *d
}
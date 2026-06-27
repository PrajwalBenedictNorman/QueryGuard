package detector


import(
	"math"
)

type CusumDetector struct{
	Target float64
	Slack float64
	HighSum float64
	LowSum float64
	Threshold float64
}

func (c *CusumDetector) CusumUpdate(value float64) (float64,bool){
	c.HighSum = math.Max(0,c.HighSum+(value- c.Target - c.Slack))
	c.LowSum =math.Max(0,c.LowSum+(c.Target-value-c.Slack))
	score :=math.Max(c.HighSum,c.LowSum)

	if score >c.Threshold{
		c.HighSum=0
		c.LowSum=0
		return score,true
	}
	return score,false
}
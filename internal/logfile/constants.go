package ionlogfile

type PeriodicRotation int

const (
	NoAutoRotate PeriodicRotation = iota
	Daily
	Weekly
	Monthly
)

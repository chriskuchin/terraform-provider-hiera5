package streamer

// A collector that doesn't accept rich data of any kind
type dataCollector struct {
	BasicCollector
}

// DataCollector returns a Collector that accepts only data, not rich-data
func DataCollector() Collector {
	c := &dataCollector{}
	c.Init()
	return c
}

func (c *dataCollector) CanDoBinary() bool {
	return false
}

func (c *dataCollector) CanDoTime() bool {
	return false
}

func (c *dataCollector) CanDoComplexKeys() bool {
	return false
}

package calucaltor

import (
	"testing"

	"finHubPipeline/structs"
)

func TestMovingAverage(t *testing.T) {
	a := New(5, nil)

	c := a.CalculateData(structs.MsgData{P: 2})
	for v := range c {
		if v < 1.999 || v > 2.001 {
			t.Fail()
		}
	}

	c = a.CalculateData(structs.MsgData{P: 4})
	c = a.CalculateData(structs.MsgData{P: 2})
	for v := range c {
		if v < 2.665 || v > 2.667 {
			t.Fail()
		}
	}
	c = a.CalculateData(structs.MsgData{P: 4})
	c = a.CalculateData(structs.MsgData{P: 2})
	for v := range c {
		if v < 2.799 || v > 2.801 {
			t.Fail()
		}
	}

	// This one will go into the first slot again
	// evicting the first value
	c = a.CalculateData(structs.MsgData{P: 10})
	for v := range c {
		if v < 4.399 || v > 4.401 {
			t.Fail()
		}
	}
}

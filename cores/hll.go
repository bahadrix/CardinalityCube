package cores

import (
	"github.com/axiomhq/hyperloglog"
)

type HLLOpts struct {
	With16Registers bool
}

type HLLCore struct {
	sketch *hyperloglog.Sketch
}

// Initialize core
var HLL CoreInitiator = func(opts interface{}) Core {
	// Get options
	var coreOpts *HLLOpts

	if opts != nil {
		coreOpts = opts.(*HLLOpts)
	} else { // Defaults
		coreOpts = &HLLOpts{
			With16Registers: false,
		}
	}

	var sketch *hyperloglog.Sketch

	if coreOpts.With16Registers {
		sketch = hyperloglog.New16()
	} else {
		sketch = hyperloglog.New()
	}

	return HLLCore{
		sketch: sketch,
	}
}

func (c HLLCore) Push(item []byte) {
	c.sketch.Insert(item)
}

func (c HLLCore) Count() uint64 {
	return c.sketch.Estimate()
}

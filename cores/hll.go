package cores

import (
	"github.com/axiomhq/hyperloglog"
)

// HLLOpts is options for HyperLogLog core
type HLLOpts struct {
	With16Registers bool // True for 16 Registers, False for 14 Registers
}

// HLLCore is a Core implementation that uses HyperLogLog.
// It is highly memory efficient with trade off against approximation.
type HLLCore struct {
	sketch *hyperloglog.Sketch
}

// HLL is a CoreInitiator implementation
func HLL(opts interface{}) Core {
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

// Push adds item to store to count.
func (c HLLCore) Push(item []byte) {
	c.sketch.Insert(item)
}

// Count returns unique item count in the store.
func (c HLLCore) Count() uint64 {
	return c.sketch.Estimate()
}

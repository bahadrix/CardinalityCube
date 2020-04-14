package cores

import (
	"fmt"
	"github.com/axiomhq/hyperloglog"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)


func TestHLLCore_Initiate(t *testing.T) {

	initializationTests := []struct{
		name          string
		initiatorCore Core
		targetCore    HLLCore
	} {
		{
			name:          "Create HLL Core with 14 Registers",
			initiatorCore: HLL(&HLLOpts{With16Registers:false}),
			targetCore:    HLLCore{sketch: hyperloglog.New14()},
		},
		{
			name:          "Create HLL Core with 16 Registers",
			initiatorCore: HLL(&HLLOpts{With16Registers:true}),
			targetCore:    HLLCore{sketch: hyperloglog.New16()},
		},
	}

	for _, it := range initializationTests {
		t.Run(it.name, func(t *testing.T) {
			if !reflect.DeepEqual(it.initiatorCore, it.targetCore) {
				t.Errorf("Initiator core %v must be %v", it.initiatorCore, it.targetCore)
			}
		})
	}

}

func TestHLLCore_Push(t *testing.T) {

	pushTests := []struct {
		totalItem  int
		repetition int
		epsilon    float64
	}{
		{
			500,
			100,
			0.005,
		},
		{
			5000,
			1000,
			0.01,
		},
		{
			500000,
			100000,
			0.01,
		},
		{
			5000000,
			1000000,
			0.01,
		},
	}

	for _, pt := range pushTests {
		testName := fmt.Sprintf("Push %d items with %d repetition", pt.totalItem, pt.repetition)
		t.Run(testName, func(t *testing.T) {

			core := HLL(nil)
			baseString := "test"

			core.Push([]byte(baseString))
			var step int = pt.totalItem / pt.repetition
			pt.repetition = pt.totalItem / step // fit for remainder


			for i := 1; i < pt.totalItem; i++ {
				var item []byte
				if i % step == 0 {
					item = []byte(baseString)
				} else {
					item = []byte(fmt.Sprintf("%s_%d", baseString, i))
				}
				core.Push(item)
			}

			realUnique := uint64(pt.totalItem - pt.repetition)
			uniqueRateActual := float64(core.Count())/float64(pt.totalItem)
			uniqueRateExpected := float64(realUnique)/float64(pt.totalItem)

			assert.InEpsilonf(t, uniqueRateExpected, uniqueRateActual, pt.epsilon, "Failure rate is too high" )
		})
	}
}
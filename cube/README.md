# Cardinality Cube

Using as a data structure example
```go
import (
	"github.com/bahadrix/cardinalitycube/cores"
	"github.com/bahadrix/cardinalitycube/cube"
)

func main() {

	hllCube := cube.CreateCube(cores.HLL, &cores.HLLOpts{
		With16Registers:false,
	})

	board := hllCube.GetBoard("sampleBoard", true)
	cell := board.GetCell("login_page", "hit", true)
	cell.Push([]byte("user_1"))

	println(cell.Count())
}
``` 
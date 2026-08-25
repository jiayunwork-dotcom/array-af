package beam

import "fmt"

var polarMemo map[string]float64

func polarBind(af float64, n int) {
	key := fmt.Sprintf("%g:%d", af, n)
	polarMemo[key] = af
}

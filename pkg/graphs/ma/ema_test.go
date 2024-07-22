package ma

import "testing"

func TestEma(t *testing.T) {
	ema := Ema{}
	ema.calculate("", 3)
}

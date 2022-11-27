package common

import (
	"fmt"
	"testing"
)

func TestBresenhamSlope1(t *testing.T) {
	fmt.Printf("%v", DrawLineOnCells(5, 5, 10, 15, 20, 20))
}

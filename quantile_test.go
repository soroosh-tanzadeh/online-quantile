package onlinequantile

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCalculateQuantilesAsExpected(t *testing.T) {
	inputs1 := []float64{0.02, 0.5, 0.74, 3.39, 0.83, 22.37, 10.15, 15.43, 38.62, 15.92, 34.60, 10.28, 1.47, 0.40, 0.05, 1.39, 0.27, 0.42, 0.09, 11.37}
	q50 := NewQuantile(0.5)

	q50.Update(inputs1[:5])

	for _, x := range inputs1[5:] {
		q50.Consume(x)
		fmt.Printf("%f %d %d %d %d %d \n", q50.GetQuantileValue(), q50.m[0].getN(), q50.m[1].getN(), q50.m[2].getN(), q50.m[3].getN(), q50.m[4].getN())
	}

	assert.Equal(t, 1, q50.GetM1().getN())
	assert.Equal(t, 6, q50.GetM2().getN())
	assert.Equal(t, 11, q50.GetM3().getN())
	assert.Equal(t, 16, q50.GetM4().getN())
	assert.Equal(t, 20, q50.GetM5().getN())

	assert.Equal(t, 4.20, math.Floor(q50.GetQuantileValue()*100)/100)
}

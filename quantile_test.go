package onlinequantile

import (
	"fmt"
	"math"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func TestShouldCalculateQuantilesAsExpected(t *testing.T) {
	series := []float64{0.02, 0.5, 0.74, 3.39, 0.83, 22.37, 10.15, 15.43, 38.62, 15.92, 34.60, 10.28, 1.47, 0.40, 0.05, 1.39, 0.27, 0.42, 0.09, 11.37}
	q50 := NewQuantile(0.5)

	q50.Update(series[:5])

	for _, x := range series[5:] {
		q50.Consume(x)
		fmt.Printf("%f %d %d %d %d %d \n", q50.GetQuantileValue(), q50.m[0].getN(), q50.m[1].getN(), q50.m[2].getN(), q50.m[3].getN(), q50.m[4].getN())
	}

	assert.Equal(t, int64(1), q50.GetM1().getN())
	assert.Equal(t, int64(6), q50.GetM2().getN())
	assert.Equal(t, int64(11), q50.GetM3().getN())
	assert.Equal(t, int64(16), q50.GetM4().getN())
	assert.Equal(t, int64(20), q50.GetM5().getN())

	assert.Equal(t, 4.20, math.Floor(q50.GetQuantileValue()*100)/100)
}

func TestExportMarkers(t *testing.T) {
	q := &Quantile{
		p: 0.5,
		m: make([]Marker, 5),
	}

	q.m[0] = &MinMarker{BaseMarker: BaseMarker{q: 1.0, n: 1, nPrime: 1.0, dPrime: 0.5}}
	q.m[1] = &MidMarker{BaseMarker: BaseMarker{q: 2.0, n: 2, nPrime: 2.0, dPrime: 0.5}}
	q.m[2] = &MidMarker{BaseMarker: BaseMarker{q: 3.0, n: 3, nPrime: 3.0, dPrime: 0.5}}
	q.m[3] = &MidMarker{BaseMarker: BaseMarker{q: 4.0, n: 4, nPrime: 4.0, dPrime: 0.5}}
	q.m[4] = &MaxMarker{BaseMarker: BaseMarker{q: 5.0, n: 5, nPrime: 5.0, dPrime: 0.5}}

	exportedMarkers := q.ExportMarkers()

	expectedMarkers := map[string]map[string]interface{}{
		"m1": {"q": 1.0, "n": int64(1), "nPrime": 1.0, "dPrime": 0.5},
		"m2": {"q": 2.0, "n": int64(2), "nPrime": 2.0, "dPrime": 0.5},
		"m3": {"q": 3.0, "n": int64(3), "nPrime": 3.0, "dPrime": 0.5},
		"m4": {"q": 4.0, "n": int64(4), "nPrime": 4.0, "dPrime": 0.5},
		"m5": {"q": 5.0, "n": int64(5), "nPrime": 5.0, "dPrime": 0.5},
	}

	if !reflect.DeepEqual(exportedMarkers, expectedMarkers) {
		t.Errorf("ExportMarkers() = %v, want %v", exportedMarkers, expectedMarkers)
	}
}

func TestFrom(t *testing.T) {
	markers := map[string]map[string]interface{}{
		"m1": {"q": 1.0, "n": int64(1), "nPrime": 1.0, "dPrime": 0.5},
		"m2": {"q": 2.0, "n": int64(2), "nPrime": 2.0, "dPrime": 0.5},
		"m3": {"q": 3.0, "n": int64(3), "nPrime": 3.0, "dPrime": 0.5},
		"m4": {"q": 4.0, "n": int64(4), "nPrime": 4.0, "dPrime": 0.5},
		"m5": {"q": 5.0, "n": int64(5), "nPrime": 5.0, "dPrime": 0.5},
	}

	q := &Quantile{
		p: 0.5,
		m: make([]Marker, 5),
	}

	q.From(markers)

	m1 := &MinMarker{BaseMarker: BaseMarker{q: 1.0, n: 1, nPrime: 1.0, dPrime: 0.5}}
	m2 := &MidMarker{BaseMarker: BaseMarker{q: 2.0, n: 2, nPrime: 2.0, dPrime: 0.5}}
	m3 := &MidMarker{BaseMarker: BaseMarker{q: 3.0, n: 3, nPrime: 3.0, dPrime: 0.5}}
	m4 := &MidMarker{BaseMarker: BaseMarker{q: 4.0, n: 4, nPrime: 4.0, dPrime: 0.5}}
	m5 := &MaxMarker{BaseMarker: BaseMarker{q: 5.0, n: 5, nPrime: 5.0, dPrime: 0.5}}

	expectedQuantile := &Quantile{
		p: 0.5,
		m: []Marker{
			m1, m2, m3, m4, m5,
		},
	}

	for i, marker := range expectedQuantile.m {
		expected := marker
		actual := q.m[i]

		// Remove Neighbors of mid markers
		if _, ok := actual.(*MidMarker); ok {
			e := actual.(*MidMarker)
			e.SetNeighbors(nil, nil)
			actual = e
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("From() markers[%d] = %v, want %v", i, actual, expected)
		}
	}
}

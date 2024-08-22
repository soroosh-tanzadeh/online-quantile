package onlinequantile

import "sort"

type Quantile struct {
	m1 Marker
	m2 Marker
	m3 Marker
	m4 Marker
	m5 Marker

	p           float64
	initialized bool
}

func (q *Quantile) GetM1() Marker {
	return q.m1
}

func (q *Quantile) GetM2() Marker {
	return q.m2
}

func (q *Quantile) GetM3() Marker {
	return q.m3
}

func (q *Quantile) GetM4() Marker {
	return q.m4
}

func (q *Quantile) GetM5() Marker {
	return q.m5
}

func (q *Quantile) update(buffer []float64) {
	i := 0

	if !q.initialized {
		q.initialize(buffer[0:4])
		i = 5
	}

	for ; i < len(buffer); i++ {
		q.Consume(buffer[i])
	}

}

func (q *Quantile) Consume(x float64) {
	q.m1.UpdatePosition(x)
	q.m2.UpdatePosition(x)
	q.m3.UpdatePosition(x)
	q.m4.UpdatePosition(x)
	q.m5.UpdatePosition(x)

	q.m2.UpdateQuantile()
	q.m3.UpdateQuantile()
	q.m4.UpdateQuantile()
}

func (q *Quantile) initialize(firstFive []float64) {
	// Sort the array
	sort.Float64s(firstFive)

	// Initialize markers
	q.m1 = &MinMarker{BaseMarker: BaseMarker{q: firstFive[0], n: 1, nPrime: 1.0, dPrime: 0.0}}

	m2 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[1], n: 2, nPrime: 1.0, dPrime: 0.0}}
	q.m2 = m2

	m3 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[2], n: 3, nPrime: 1.0, dPrime: 0.0}}
	q.m3 = m3

	m4 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[3], n: 4, nPrime: 1.0, dPrime: 0.0}}
	q.m4 = m4

	q.m5 = &MaxMarker{BaseMarker: BaseMarker{q: firstFive[4], n: 5, nPrime: 1.0, dPrime: 0.0}}

	// Set neighbors
	m2.SetNeighbors(q.m1, q.m3)
	m3.SetNeighbors(q.m2, q.m4)
	m4.SetNeighbors(q.m3, q.m5)

	q.initialized = true
}

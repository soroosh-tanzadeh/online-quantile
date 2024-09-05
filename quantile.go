package onlinequantile

import "sort"

type Quantile struct {
	m []Marker

	p           float64
	initialized bool
}

func NewQuantile(p float64) *Quantile {
	return &Quantile{
		p: p,
		m: make([]Marker, 5),
	}
}

func (q *Quantile) GetM1() Marker {
	return q.m[0]
}

func (q *Quantile) GetM2() Marker {
	return q.m[1]
}

func (q *Quantile) GetM3() Marker {
	return q.m[2]
}

func (q *Quantile) GetM4() Marker {
	return q.m[3]
}

func (q *Quantile) GetM5() Marker {
	return q.m[4]
}

func (q *Quantile) Update(buffer []float64) {
	i := 0

	if !q.initialized {
		q.initialize(buffer[0:5])
		i = 5
	}

	for ; i < len(buffer); i++ {
		q.Consume(buffer[i])
	}

}

func (q *Quantile) GetQuantileValue() float64 {
	return q.m[2].GetValue()
}

func (q *Quantile) Consume(x float64) {
	var k int
	if x < q.m[0].GetValue() {
		q.m[0].SetValue(x)
		k = 0
	}
	if x >= q.m[0].GetValue() && x < q.m[1].GetValue() {
		k = 0
	}
	if x >= q.m[1].GetValue() && x < q.m[2].GetValue() {
		k = 1
	}
	if x >= q.m[4].GetValue() && x < q.m[3].GetValue() {
		k = 2
	}
	if x >= q.m[3].GetValue() && x < q.m[4].GetValue() {
		k = 3
	}
	if x > q.m[4].GetValue() {
		k = 3
		q.m[4].SetValue(x)
	}

	for i := 0; i < 5; i++ {
		if i >= k+1 {
			q.m[i].IncrementPosition()
		}
		q.m[i].IncrementDesiredPosition()
	}

	q.m[1].UpdateQuantile()
	q.m[2].UpdateQuantile()
	q.m[3].UpdateQuantile()
}

func (q *Quantile) initialize(firstFive []float64) {
	// Sort the array
	sort.Float64s(firstFive)

	// Initialize markers
	q.m[0] = &MinMarker{BaseMarker: BaseMarker{q: firstFive[0], n: 1, nPrime: 1.0, dPrime: 0.0}}

	m2 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[1], n: 2, nPrime: 1.0 + (2 * q.p), dPrime: q.p / 2.0}}
	q.m[1] = m2

	m3 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[2], n: 3, nPrime: 1.0 + (4 * q.p), dPrime: q.p}}
	q.m[2] = m3

	m4 := &MidMarker{BaseMarker: BaseMarker{q: firstFive[3], n: 4, nPrime: 3.0 + (2 * q.p), dPrime: (1.0 + q.p) / 2.0}}
	q.m[3] = m4

	q.m[4] = &MaxMarker{BaseMarker: BaseMarker{q: firstFive[4], n: 5, nPrime: 5, dPrime: 1.0}}

	// Set neighbors
	m2.SetNeighbors(q.m[0], q.m[2])
	m3.SetNeighbors(q.m[1], q.m[3])
	m4.SetNeighbors(q.m[2], q.m[4])

	q.initialized = true
}

func (q *Quantile) From(markers map[string]map[string]any) {

	// Initialize markers using the provided markers map
	q.m[0] = &MinMarker{
		BaseMarker: BaseMarker{
			q:      markers["m1"]["q"].(float64),
			n:      markers["m1"]["n"].(int64),
			nPrime: markers["m1"]["nPrime"].(float64),
			dPrime: markers["m1"]["dPrime"].(float64),
		},
	}

	m2 := &MidMarker{
		BaseMarker: BaseMarker{
			q:      markers["m2"]["q"].(float64),
			n:      markers["m2"]["n"].(int64),
			nPrime: markers["m2"]["nPrime"].(float64),
			dPrime: markers["m2"]["dPrime"].(float64),
		},
	}
	q.m[1] = m2

	m3 := &MidMarker{
		BaseMarker: BaseMarker{
			q:      markers["m3"]["q"].(float64),
			n:      markers["m3"]["n"].(int64),
			nPrime: markers["m3"]["nPrime"].(float64),
			dPrime: markers["m3"]["dPrime"].(float64),
		},
	}
	q.m[2] = m3

	m4 := &MidMarker{
		BaseMarker: BaseMarker{
			q:      markers["m4"]["q"].(float64),
			n:      markers["m4"]["n"].(int64),
			nPrime: markers["m4"]["nPrime"].(float64),
			dPrime: markers["m4"]["dPrime"].(float64),
		},
	}
	q.m[3] = m4

	q.m[4] = &MaxMarker{
		BaseMarker: BaseMarker{
			q:      markers["m5"]["q"].(float64),
			n:      markers["m5"]["n"].(int64),
			nPrime: markers["m5"]["nPrime"].(float64),
			dPrime: markers["m5"]["dPrime"].(float64),
		},
	}

	// Set neighbors
	m2.SetNeighbors(q.m[0], q.m[2])
	m3.SetNeighbors(q.m[1], q.m[3])
	m4.SetNeighbors(q.m[2], q.m[4])

	q.initialized = true
}

func (q *Quantile) ExportMarkers() map[string]map[string]any {
	m1 := q.GetM1().(*MinMarker)
	m2 := q.GetM2().(*MidMarker)
	m3 := q.GetM3().(*MidMarker)
	m4 := q.GetM4().(*MidMarker)
	m5 := q.GetM5().(*MaxMarker)

	return map[string]map[string]any{
		"m1": {
			"q":      m1.q,
			"n":      m1.n,
			"dPrime": m1.dPrime,
			"nPrime": m1.nPrime,
		},
		"m2": {
			"q":      m2.q,
			"n":      m2.n,
			"dPrime": m2.dPrime,
			"nPrime": m2.nPrime,
		},
		"m3": {
			"q":      m3.q,
			"n":      m3.n,
			"dPrime": m3.dPrime,
			"nPrime": m3.nPrime,
		},
		"m4": {
			"q":      m4.q,
			"n":      m4.n,
			"dPrime": m4.dPrime,
			"nPrime": m4.nPrime,
		},
		"m5": {
			"q":      m5.q,
			"n":      m5.n,
			"dPrime": m5.dPrime,
			"nPrime": m5.nPrime,
		},
	}
}

package onlinequantile

type Marker interface {
	UpdatePosition(float64)
	UpdateQuantile()
	GetValue() float64

	getN() int
}

type BaseMarker struct {
	q      float64
	n      int
	nPrime float64
	dPrime float64
}

func (m *BaseMarker) GetValue() float64 {
	return m.q
}

func (m *BaseMarker) getN() int {
	return m.n
}

func (m *BaseMarker) updateDesiredPosition() {
	m.nPrime += m.dPrime
}

func (m *BaseMarker) UpdateQuantile() {
	// Nothing
}

func (m *BaseMarker) UpdatePosition(float64) {
	// Nothing
}

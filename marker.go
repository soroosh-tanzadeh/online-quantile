package onlinequantile

type Marker interface {
	IncrementPosition()
	IncrementDesiredPosition()
	UpdateQuantile()
	GetValue() float64
	SetValue(x float64)
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

func (m *BaseMarker) SetValue(x float64) {
	m.q = x
}

func (m *BaseMarker) updateDesiredPosition() {
	m.nPrime += m.dPrime
}

func (m *BaseMarker) UpdateQuantile() {
	// Nothing
}

func (m *BaseMarker) IncrementPosition() {
	m.n++
}

func (m *BaseMarker) IncrementDesiredPosition() {
	m.nPrime = m.nPrime + m.dPrime
}

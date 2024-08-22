package onlinequantile

type MaxMarker struct {
	BaseMarker
}

func (m *MaxMarker) UpdatePosition(v float64) {
	m.n++
	m.updateDesiredPosition()

	if v > m.q {
		m.q = v
	}
}

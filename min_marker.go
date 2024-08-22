package onlinequantile

type MinMarker struct {
	BaseMarker
}

func (m *MinMarker) UpdatePosition(v float64) {
	if v < m.q {
		m.q = v
	}
}

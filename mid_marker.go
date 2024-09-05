package onlinequantile

type MidMarker struct {
	BaseMarker
	lNeighbor Marker
	rNeighbor Marker
}

func (m *MidMarker) UpdatePosition(value float64) {
	if value < m.q {
		m.n++
	}

	m.updateDesiredPosition()
}

func (m *MidMarker) SetNeighbors(left, right Marker) {
	m.lNeighbor = left
	m.rNeighbor = right
}

func (m *MidMarker) UpdateQuantile() {
	offsetFromDesired := m.nPrime - float64(m.n)
	offsetFromRNeighbor := m.rNeighbor.getN() - m.n
	offsetFromLNeighbor := m.lNeighbor.getN() - m.n

	var displacemment int64 = 1
	if offsetFromDesired < 0.0 {
		displacemment = -1
	}

	if (offsetFromDesired >= 1.0 && offsetFromRNeighbor > 1) || (offsetFromDesired <= -1 && offsetFromLNeighbor < -1) {
		qTmp := m.pSquared(displacemment)
		if m.lNeighbor.GetValue() >= qTmp || qTmp >= m.rNeighbor.GetValue() {
			qTmp = m.linear(displacemment)
		}

		m.q = qTmp
		m.n += displacemment
	}

}

func (m *MidMarker) pSquared(displacement int64) float64 {
	neighborSpan := m.rNeighbor.getN() - m.lNeighbor.getN()
	lNeighborOffset := m.n - m.lNeighbor.getN()
	rNeighborOffset := m.rNeighbor.getN() - m.n

	qDifRNeighbor := m.rNeighbor.GetValue() - m.q
	qDifLNeighbor := m.q - m.lNeighbor.GetValue()

	return (m.q + float64(displacement)/float64(neighborSpan)) * ((float64(lNeighborOffset+displacement) * (qDifRNeighbor / float64(rNeighborOffset))) +
		(float64(rNeighborOffset-displacement) * (qDifLNeighbor / float64(lNeighborOffset))))
}

func (m *MidMarker) linear(displacement int64) float64 {
	neighborQ := 0.0
	if displacement < 0 {
		neighborQ = m.lNeighbor.GetValue()
	} else {
		neighborQ = m.rNeighbor.GetValue()
	}

	var neighborOffset int64 = 0
	if displacement < 0 {
		neighborOffset = m.lNeighbor.getN()
	} else {
		neighborOffset = m.rNeighbor.getN()
	}

	qDiff := neighborQ - m.q
	offsetDiff := neighborOffset - m.n

	return m.q + float64(displacement)*qDiff/float64(offsetDiff)
}

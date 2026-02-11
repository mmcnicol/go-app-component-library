// pkg/components/chart/streaming_data.go
package chart

type StreamingData struct {
	Points   []float64
	Capacity int
}

func (s *StreamingData) Push(val float64) {
	if len(s.Points) >= s.Capacity {
		s.Points = s.Points[1:] // Remove oldest
	}
	s.Points = append(s.Points, val)
}

// pkg/components/chart/constants.go
package chart

import "math"

func generateRandomData(count int) []DataPoint {
	var points []DataPoint
	for i := 0; i < count; i++ {
		points = append(points, DataPoint{
			X: float64(i),
			Y: float64(i) + (math.Sin(float64(i)*0.1) * 10),
		})
	}
	return points
}

// pkg/components/chart/constants.go
package chart

var userActivity = [][]float64{
	{10, 20, 30, 40, 50, 60, 70, 80},
	{15, 25, 35, 45, 55, 65, 75, 85},
	{5, 15, 25, 35, 45, 55, 65, 75},
}

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

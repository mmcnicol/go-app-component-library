// pkg/components/chart/internal/math.go
package internal

import "github.com/mmcnicol/go-app-component-library/pkg/components/chart"

func CalculateLinearRegression(data []Point) (m, b float64) {
	n := float64(len(data))
	if n < 2 {
		return 0, 0
	}

	var sumX, sumY, sumXY, sumXX float64
	for _, p := range data {
		sumX += p.X
		sumY += p.Y
		sumXY += p.X * p.Y
		sumXX += p.X * p.X
	}

	denominator := (n * sumXX) - (sumX * sumX)
	if denominator == 0 {
		return 0, 0
	}

	m = (n*sumXY - sumX*sumY) / denominator
	b = (sumY - m*sumX) / n
	return m, b
}

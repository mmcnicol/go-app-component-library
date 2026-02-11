// pkg/components/chart/utils.go
package chart

func CalculateBoxStats(data []float64) BoxPlotStats {
	// Sort the data first
	sort.Float64s(data)
	n := len(data)
	
	return BoxPlotStats{
		Min:    data[0],
		Q1:     data[n/4],
		Median: data[n/2],
		Q3:     data[3*n/4],
		Max:    data[n-1],
	}
}

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

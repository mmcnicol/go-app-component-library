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

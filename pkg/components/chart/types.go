// pkg/components/chart/types.go
package chart

type Point struct {
	X, Y float64
}

type BoxPlotStats struct {
	Min, Q1, Median, Q3, Max float64
}

// pkg/components/chart/box_plot_chart.go
package chart

type BoxPlotChart struct {
    BaseChart
    showOutliers bool
    showMean     bool
    whiskerType  WhiskerType // "tukey", "minmax", "percentile"
}

func (bpc *BoxPlotChart) calculateStatistics(data [][]float64) []BoxPlotStats {
    stats := make([]BoxPlotStats, len(data))
    
    for i, dataset := range data {
        // Sort data
        sorted := make([]float64, len(dataset))
        copy(sorted, dataset)
        sort.Float64s(sorted)
        
        // Calculate quartiles
        q1 := bpc.calculatePercentile(sorted, 25)
        median := bpc.calculatePercentile(sorted, 50)
        q3 := bpc.calculatePercentile(sorted, 75)
        
        // Calculate IQR
        iqr := q3 - q1
        
        // Determine whiskers
        var lowerWhisker, upperWhisker float64
        switch bpc.whiskerType {
        case WhiskerTypeTukey:
            lowerWhisker = q1 - 1.5*iqr
            upperWhisker = q3 + 1.5*iqr
        case WhiskerTypeMinMax:
            lowerWhisker = sorted[0]
            upperWhisker = sorted[len(sorted)-1]
        case WhiskerTypePercentile:
            lowerWhisker = bpc.calculatePercentile(sorted, 5)
            upperWhisker = bpc.calculatePercentile(sorted, 95)
        }
        
        // Identify outliers
        var outliers []float64
        if bpc.showOutliers {
            for _, value := range sorted {
                if value < lowerWhisker || value > upperWhisker {
                    outliers = append(outliers, value)
                }
            }
        }
        
        stats[i] = BoxPlotStats{
            Min:          sorted[0],
            Q1:           q1,
            Median:       median,
            Q3:           q3,
            Max:          sorted[len(sorted)-1],
            LowerWhisker: lowerWhisker,
            UpperWhisker: upperWhisker,
            Outliers:     outliers,
            Mean:         bpc.calculateMean(dataset),
        }
    }
    
    return stats
}

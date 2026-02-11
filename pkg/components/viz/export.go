// pkg/components/viz/export.go
package viz

import (
    "fmt"
    "strings"
    "time"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ChartExporter handles exporting charts to various formats
type ChartExporter struct {
    chart *Chart
}

func NewChartExporter(chart *Chart) *ChartExporter {
    return &ChartExporter{chart: chart}
}

func (e *ChartExporter) Export(format ExportFormat, options map[string]interface{}) error {
    switch format {
    case ExportFormatPNG:
        return e.exportPNG()
    case ExportFormatSVG:
        return e.exportSVG()
    case ExportFormatPDF:
        return e.exportPDF()
    case ExportFormatCSV:
        return e.exportCSV()
    case ExportFormatJSON:
        return e.exportJSON()
    default:
        return fmt.Errorf("unsupported export format: %v", format)
    }
}

func (e *ChartExporter) exportPNG() error {
    // Placeholder for PNG export
    app.Log("PNG export not yet implemented")
    return nil
}

func (e *ChartExporter) exportSVG() error {
    app.Log("SVG export not yet implemented")
    return nil
}

func (e *ChartExporter) exportPDF() error {
    app.Log("PDF export not yet implemented")
    return nil
}

func (e *ChartExporter) exportCSV() error {
    var csv strings.Builder
    
    // Write headers
    csv.WriteString("Label")
    for _, series := range e.chart.spec.Data.Series {
        csv.WriteString("," + series.Label)
    }
    csv.WriteString("\n")
    
    // Write data
    for i, label := range e.chart.spec.Data.Labels {
        csv.WriteString(label)
        for _, series := range e.chart.spec.Data.Series {
            if i < len(series.Points) {
                csv.WriteString(fmt.Sprintf(",%f", series.Points[i].Y))
            } else {
                csv.WriteString(",")
            }
        }
        csv.WriteString("\n")
    }
    
    // Trigger download
    filename := fmt.Sprintf("chart_%s.csv", time.Now().Format("20060102_150405"))
    app.Window().Call("downloadText", csv.String(), filename)
    
    return nil
}

func (e *ChartExporter) exportJSON() error {
    app.Log("JSON export not yet implemented")
    return nil
}

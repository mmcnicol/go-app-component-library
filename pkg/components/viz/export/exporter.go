// pkg/components/viz/export/exporter.go
package export

type Format string
const (
    PNG Format = "png"
    SVG Format = "svg"
    PDF Format = "pdf"
    CSV Format = "csv"
    JSON Format = "json"
)

type Exporter struct {
    chart *viz.Chart
}

func (e *Exporter) PNG() error {
    // Export as PNG image
}

func (e *Exporter) CSV() error {
    // Export data as CSV
}

func (e *Exporter) SVG() error {
    // Export as SVG vector graphics
}

// pkg/components/vizzy/render_canvas_bar_chart.go
package canvas

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type BarChartRenderer struct {
    ctx     app.Value
    width   int
    height  int
    margins Margins
}

func (r *BarChartRenderer) Render(chart *Chart) error {
    // Pure rendering logic - no chart logic
    r.calculateScales(chart)
    r.drawBackground(chart)
    r.drawAxes(chart)
    r.drawBars(chart)
    r.drawLabels(chart)
    r.drawTitle(chart)
    return nil
}

// Register this renderer
func init() {
    RegisterEngine(ChartTypeBar, EngineTypeCanvas, func() ChartEngine {
        return &CanvasEngine{
            renderers: map[ChartType]Renderer{
                ChartTypeBar: &BarChartRenderer{},
                ChartTypeLine: &LineChartRenderer{},
                // ... etc
            },
        }
    })
}

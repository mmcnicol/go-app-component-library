// pkg/components/vizzy/component_chart.go
package viz

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ChartComponent is the app.UI wrapper
type ChartComponent struct {
    app.Compo
    
    chartModel *Chart
    engine     ChartEngine
    canvasID   string
}

func NewChartComponent(model *Chart) *ChartComponent {
    return &ChartComponent{
        chartModel: model,
        canvasID:   generateCanvasID(),
    }
}

func (c *ChartComponent) OnMount(ctx app.Context) {
    // Get canvas element
    canvasElem := app.Window().GetElementByID(c.canvasID)
    if !canvasElem.Truthy() {
        app.Log("Canvas not found")
        return
    }
    
    // Initialize engine based on chart type
    engine, err := GetEngine(c.chartModel.Type, c.chartModel.Engine)
    if err != nil {
        app.Log("Failed to get engine:", err)
        return
    }
    
    c.engine = engine
    err = c.engine.Initialize(c.chartModel.Width, c.chartModel.Height, canvasElem)
    if err != nil {
        app.Log("Failed to initialize engine:", err)
        return
    }
    
    // Render
    c.engine.Render(c.chartModel)
}

func (c *ChartComponent) Render() app.UI {
    return app.Canvas().
        ID(c.canvasID).
        Width(c.chartModel.Width).
        Height(c.chartModel.Height).
        Style("width", "100%").
        Style("height", "100%").
        Style("display", "block")
}

// pkg/components/chart/chart_engine_wrapper.go
package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ChartEngineWrapper wraps CanvasRenderer to implement ChartEngine interface
type ChartEngineWrapper struct {
    renderer *CanvasRenderer
}

func NewChartEngineWrapper(containerID string) (*ChartEngineWrapper, error) {
    renderer, err := NewCanvasRenderer(containerID)
    if err != nil {
        return nil, err
    }
    return &ChartEngineWrapper{
        renderer: renderer,
    }, nil
}

// Render implements ChartEngine.Render
func (cw *ChartEngineWrapper) Render(chart ChartSpec) error {
    return cw.renderer.RenderChart(chart)
}

// Update implements ChartEngine.Update
func (cw *ChartEngineWrapper) Update(data ChartData) error {
    return cw.renderer.Update(data)
}

// Destroy implements ChartEngine.Destroy
func (cw *ChartEngineWrapper) Destroy() error {
    return cw.renderer.Destroy()
}

// GetCanvas implements ChartEngine.GetCanvas
func (cw *ChartEngineWrapper) GetCanvas() app.UI {
    return cw.renderer.GetCanvas()
}

// pkg/components/viz/engines/engine.go
package engines

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/yourlib/pkg/components/viz"
)

// Engine defines the chart rendering interface
type Engine interface {
    // Lifecycle
    Init(canvas app.HTMLCanvas) error
    Render(spec *viz.Spec) error
    Update(data viz.DataSet) error
    Resize(width, height int) error
    Destroy() error
    
    // Interactive
    HitTest(x, y float64) ([]viz.Point, int, error)
    GetCanvas() app.HTMLCanvas
    
    // Performance
    SetMaxPoints(n int)
    GetMetrics() Metrics
}

// CanvasEngine - Pure Go canvas implementation
type CanvasEngine struct {
    // ... implementation
}

// WebGLEngine - High-performance WebGL for large datasets
type WebGLEngine struct {
    // ... implementation
}

// SVGRenderer - For export quality
type SVGRenderer struct {
    // ... implementation
}

// AutoEngine - Automatically selects best engine
func AutoEngine(spec *viz.Spec) Engine {
    pointCount := spec.Data.PointCount()
    
    switch {
    case pointCount > 50000:
        return NewWebGLEngine()
    case pointCount > 1000:
        return NewCanvasEngine()
    default:
        return NewCanvasEngine()
    }
}

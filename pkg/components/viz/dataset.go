// pkg/components/viz/dataset.go
package viz

import (
    "fmt"
)

// Stroke defines line styling
type Stroke struct {
    Width float64 `json:"width"`
    Color string  `json:"color"`
    Dash  []float64 `json:"dash,omitempty"`
}

// PointStyle defines point marker style
type PointStyle string

const (
    PointStyleCircle   PointStyle = "circle"
    PointStyleSquare   PointStyle = "square"
    PointStyleTriangle PointStyle = "triangle"
    PointStyleCross    PointStyle = "cross"
    PointStyleDiamond  PointStyle = "diamond"
)

// Node represents a node in hierarchical data structures
type Node struct {
    Name     string      `json:"name"`
    Value    float64     `json:"value,omitempty"`
    Children []*Node     `json:"children,omitempty"`
    Data     interface{} `json:"data,omitempty"`
}

// DataSet represents a complete dataset for visualization
type DataSet struct {
    // Labels for categories (x-axis, slices, etc)
    Labels    []string `json:"labels"`
    
    // Multiple series support
    Series    []Series `json:"series"`
    
    // 2D data for heatmaps, matrices
    Matrix    [][]float64 `json:"matrix,omitempty"`
    
    // Hierarchical data for tree maps
    Hierarchy *Node `json:"hierarchy,omitempty"`
    
    // Metadata
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Series represents a single data series
type Series struct {
    // Identity
    Label     string      `json:"label"`
    
    // Data points
    Points    []Point     `json:"points"`
    
    // Appearance (can override theme)
    Color     string      `json:"color,omitempty"`
    Fill      bool        `json:"fill"`
    Stroke    Stroke      `json:"stroke,omitempty"`
    
    // Point styling
    PointSize int         `json:"pointSize"`
    PointStyle PointStyle `json:"pointStyle"`
    
    // Additional
    Stack     string      `json:"stack,omitempty"` // for stacked charts
    YAxisID   string      `json:"yAxisID,omitempty"` // for multi-axis

    // Line smoothing (0-1)
    Tension   float64     `json:"tension,omitempty"`
}

// Point represents a single data point
type Point struct {
    X     float64     `json:"x"`
    Y     float64     `json:"y"`
    R     float64     `json:"r,omitempty"` // radius for bubble charts
    Label string      `json:"label,omitempty"`
    
    // Raw value for pie charts
    Value float64     `json:"value,omitempty"`
}

// Convenience constructors
func XYPoints(x, y []float64) []Point {
    points := make([]Point, len(x))
    for i := range x {
        points[i] = Point{X: x[i], Y: y[i]}
    }
    return points
}

func Values(values []float64) []Point {
    points := make([]Point, len(values))
    for i, v := range values {
        points[i] = Point{Y: v, Label: fmt.Sprintf("%d", i+1)}
    }
    return points
}

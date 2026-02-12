// pkg/components/viz/core_types.go
package core

// ChartType defines available chart types
type ChartType string

const (
    ChartTypeLine        ChartType = "line"
    ChartTypeBar         ChartType = "bar"
    ChartTypeGroupedBar  ChartType = "grouped-bar"
    ChartTypeStackedBar  ChartType = "stacked-bar"
    ChartTypePie         ChartType = "pie"
    ChartTypeDonut       ChartType = "donut"
    // ... all other chart types
)

// EngineType defines rendering engine
type EngineType string

const (
    EngineTypeAuto   EngineType = "auto"
    EngineTypeCanvas EngineType = "canvas"
    EngineTypeWebGL  EngineType = "webgl"
    EngineTypeSVG    EngineType = "svg"
)

// ... move ALL other constants (TooltipMode, RegressionType, etc.)

// Point represents a single data point
type Point struct {
    X     float64 `json:"x"`
    Y     float64 `json:"y"`
    R     float64 `json:"r,omitempty"`
    Label string  `json:"label,omitempty"`
    Value float64 `json:"value,omitempty"`
}

// Series represents a single data series
type Series struct {
    Label      string    `json:"label"`
    Points     []Point   `json:"points"`
    Color      string    `json:"color,omitempty"`
    Fill       bool      `json:"fill"`
    Stroke     Stroke    `json:"stroke,omitempty"`
    PointSize  int       `json:"pointSize"`
    PointStyle PointStyle `json:"pointStyle"`
    Stack      string    `json:"stack,omitempty"`
    YAxisID    string    `json:"yAxisID,omitempty"`
    Tension    float64   `json:"tension,omitempty"`
}

// DataSet represents a complete dataset
type DataSet struct {
    Labels    []string               `json:"labels"`
    Series    []Series              `json:"series"`
    Matrix    [][]float64           `json:"matrix,omitempty"`
    Hierarchy *Node                 `json:"hierarchy,omitempty"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ... move ALL other shared structs (Theme, Configs, etc.)

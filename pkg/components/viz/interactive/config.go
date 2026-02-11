// pkg/components/viz/interactive/config.go
package viz

// InteractiveConfig defines interactivity options
type InteractiveConfig struct {
    Enabled     bool
    
    // Tooltip
    Tooltip     TooltipConfig
    
    // Zoom & Pan
    Zoom        ZoomConfig
    Pan         PanConfig
    
    // Selection
    Selection   SelectionConfig
    
    // Events
    OnClick     func(ctx app.Context, e Event, points []Point)
    OnHover     func(ctx app.Context, e Event, point *Point)
    OnZoom      func(ctx app.Context, range AxisRange)
}

// TooltipConfig defines tooltip appearance and behavior
type TooltipConfig struct {
    Enabled     bool
    Mode        TooltipMode // "single", "all", "nearest"
    Intersect   bool        // Only show when mouse intersects point
    Format      func(point Point, series Series) string
    Background  string
    TextColor   string
    BorderColor string
    BorderWidth int
    Padding     int
}

// TooltipMode constants
type TooltipMode string
const (
    TooltipModeSingle   TooltipMode = "single"
    TooltipModeAll      TooltipMode = "all"
    TooltipModeNearest  TooltipMode = "nearest"
)

// pkg/components/viz/api.go
package viz

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Chart is the main component
type Chart struct {
    app.Compo
    spec *Spec
}

// Spec defines the complete chart specification
type Spec struct {
    // Identity
    ID      string
    Title   string
    
    // Content
    Type    ChartType
    Data    DataSet
    
    // Appearance
    Theme   Theme
    Width   int
    Height  int
    
    // Features
    Interactive InteractiveConfig
    Accessible  AccessibilityConfig
    Animated    bool
    
    // Performance
    Engine      EngineType // "auto", "canvas", "webgl"
    MaxPoints   int
}

// New creates a new chart with the given spec
func New(spec Spec) *Chart {
    // Apply defaults
    if spec.Engine == "" {
        spec.Engine = EngineTypeAuto
    }
    if spec.Theme == nil {
        spec.Theme = DefaultTheme()
    }
    
    return &Chart{
        spec: &spec,
    }
}

// Chainable methods for fluent API
func (c *Chart) WithTitle(title string) *Chart {
    c.spec.Title = title
    return c
}

func (c *Chart) WithData(data DataSet) *Chart {
    c.spec.Data = data
    return c
}

func (c *Chart) WithTheme(theme Theme) *Chart {
    c.spec.Theme = theme
    return c
}

func (c *Chart) Interactive() *Chart {
    c.spec.Interactive.Enabled = true
    return c
}

func (c *Chart) Accessible() *Chart {
    c.spec.Accessible.Enabled = true
    return c
}

func (c *Chart) Animated() *Chart {
    c.spec.Animated = true
    return c
}

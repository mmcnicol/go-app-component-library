// pkg/components/viz/config.go
package viz

// This file should only contain extensions and additional helper functions
// for the InteractiveConfig, not redeclarations

// ApplyDefaults sets sensible defaults for InteractiveConfig
func (c *InteractiveConfig) ApplyDefaults() {
    if c.Tooltip.IntersectDistance == 0 {
        c.Tooltip.IntersectDistance = 20
    }
    if c.Zoom.Factor == 0 {
        c.Zoom.Factor = 1.1
    }
}

// WithTooltipFormat sets a custom tooltip formatter
func (c *InteractiveConfig) WithTooltipFormat(format func(point Point, series Series) string) *InteractiveConfig {
    c.Tooltip.Format = format
    return c
}

// WithTooltipStyle sets custom tooltip styling
func (c *InteractiveConfig) WithTooltipStyle(background, textColor, borderColor string, borderWidth, padding int) *InteractiveConfig {
    c.Tooltip.Background = background
    c.Tooltip.TextColor = textColor
    c.Tooltip.BorderColor = borderColor
    c.Tooltip.BorderWidth = borderWidth
    c.Tooltip.Padding = padding
    return c
}

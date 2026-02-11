//go:build dev

package chart

import (
	"math"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    rows, cols := 20, 20
    heatmapData := make([][]float64, rows)
    for r := 0; r < rows; r++ {
        heatmapData[r] = make([]float64, cols)
        for l := 0; l < cols; l++ {
            val := (math.Sin(float64(r)*0.3) + math.Cos(float64(l)*0.3) + 2) / 4 * 100
            heatmapData[r][l] = val
        }
    }

    storybook.Register("Charts", "Heatmap",
        map[string]*storybook.Control{
            "Title":         storybook.NewTextControl("Resource Density"),
            "Color Scheme":  storybook.NewSelectControl([]string{"Blue-Green", "Inferno", "Grayscale"}, "Blue-Green"),
            "Cell Opacity":  storybook.NewRangeControl(0, 100, 10, 100),
        },
        func(controls map[string]*storybook.Control) app.UI {
            return New(nil,
                WithTitle(controls["Title"].Value.(string)),
                func(c *ChartConfig) {
                    // Clear other chart types
                    c.BoxData = nil
                    c.PieData = nil
                    c.IsStream = false
                    
                    // Set heatmap data
                    c.HeatmapMatrix = heatmapData
                    c.ColorScheme = controls["Color Scheme"].Value.(string)
                    c.Opacity = float64(controls["Cell Opacity"].Value.(int)) / 100.0
                },
            )
        },
    )
}

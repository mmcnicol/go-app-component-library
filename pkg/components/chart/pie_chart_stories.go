//go:build dev

package chart

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    // Preset datasets
    datasets := map[string][]float64{
        "Browser Market Share": {65, 18, 9, 8},
        "Device Usage":         {40, 30, 20, 10},
        "Equal Split":          {25, 25, 25, 25},
    }

    storybook.Register("Charts", "Pie Chart",
        map[string]*storybook.Control{
            "Dataset":      storybook.NewSelectControl([]string{"Browser Market Share", "Device Usage", "Equal Split"}, "Browser Market Share"),
            "Show Labels":  storybook.NewBoolControl(true),
            "Inner Radius": storybook.NewRangeControl(0, 80, 5, 0), // 0 = Pie, >0 = Donut
        },
        func(controls map[string]*storybook.Control) app.UI {
            selectedKey := controls["Dataset"].Value.(string)
            data := datasets[selectedKey]

            return New(nil,
                WithTitle(selectedKey),
                func(c *ChartConfig) {
                    c.PieData = data
                    c.InnerRadiusRatio = float64(controls["Inner Radius"].Value.(int)) / 100.0
                },
            )
        },
    )
}

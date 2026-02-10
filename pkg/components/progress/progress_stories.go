//go:build dev
// pkg/components/progress/progress_stories.go
package progress

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    storybook.Register("Misc", "Progress", 
        map[string]*storybook.Control{
            "Value": {Label: "Value", Type: storybook.ControlNumber, Value: 75},
        },
        func(controls map[string]*storybook.Control) app.UI {
            val := float64(controls["Value"].Value.(int))
            
            return app.Div().Style("padding", "20px").Body(
                &Progress{
                    Value: val,
                },
            )
        },
    )

}

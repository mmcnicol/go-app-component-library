//go:build dev
// pkg/components/label/label_stories.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	storybook.Register("Form", "Label", 
		map[string]*storybook.Control{
			"Text":     {Label: "Label Text", Type: storybook.ControlText, Value: "First Name"},
			"Required": {Label: "Required", Type: storybook.ControlBool, Value: true},
		},
		func(controls map[string]*storybook.Control) app.UI {
			return app.Div().Style("padding", "40px").Body(
				&Label{
					Text:     controls["Text"].Value.(string),
					Required: controls["Required"].Value.(bool),
				},
				// Adding a dummy input to show how it looks in a layout
				app.Input().
					Style("display", "block").
					Style("width", "200px").
					Style("padding", "8px").
					//Placeholder("Type here..."),
			)
		},
	)
}

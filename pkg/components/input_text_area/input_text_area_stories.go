//go:build dev
// pkg/components/input_text_area/input_text_area_stories.go
package input_text_area

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	storybook.Register("Form", "InputTextArea", 
		map[string]*storybook.Control{
			"Value":       {Label: "Value", Type: storybook.ControlText, Value: ""}, 
			"Placeholder": {Label: "Placeholder", Type: storybook.ControlText, Value: "Enter your message"}, 
			"Rows":        {Label: "Rows", Type: storybook.ControlNumber, Value: 5},
			"Cols":        {Label: "Cols", Type: storybook.ControlNumber, Value: 30},
			"Disabled":    {Label: "Disabled", Type: storybook.ControlBool, Value: false},
		},
		func(controls map[string]*storybook.Control) app.UI {
			val := controls["Value"].Value.(string)
			placeholder := controls["Placeholder"].Value.(string)
			rows := controls["Rows"].Value.(int)
			cols := controls["Cols"].Value.(int)
			isDisabled := controls["Disabled"].Value.(bool)

			return &InputTextArea{
				Value:       val,
				Placeholder: placeholder,
				Rows:        rows,
				Cols:        cols,
				Disabled:    isDisabled,
				OnInput: func(ctx app.Context, e app.Event) {
					// Sync the value back to the storybook control for real-time debugging
					newVal := ctx.JSSrc().Get("value").String()
					controls["Value"].Value = newVal
					ctx.Update()
				},
			}
		},
	)

}

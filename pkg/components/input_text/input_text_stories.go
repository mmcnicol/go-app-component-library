//go:build dev
// pkg/components/input_text/input_text_stories.go
package input_text

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	/*
	inputTextDefault := &InputText{
		Label: "Label text.",
	}

	storybook.Register("Input Text", "Default", func() app.UI {
		return inputTextDefault
	})

	inputTextPlaceholder := &InputText{
		Label: "Label text.",
		Placeholder: "sample placeholder text.",
	}

	storybook.Register("Input Text", "Placeholder", func() app.UI {
		return inputTextPlaceholder
	})

	inputTextReadOnly := &InputText{
		Label: "Label text.",
		Value: "this is a test",
		Disabled: true,
	}

	storybook.Register("Input Text", "ReadOnly", func() app.UI {
		return inputTextReadOnly
	})
	*/

	storybook.Register("Built In", "Input Text", 
		map[string]*storybook.Control{
			"Label": {Label: "Label Text", Type: storybook.ControlText, Value: "Label Text."},
			"Value": {Label: "Value", Type: storybook.ControlText, Value: ""}, 
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Placeholder": {Label: "Placeholder", Type: storybook.ControlText, Value: "First Name"}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			labelText := controls["Label Text"].Value.(string)
			valueString := controls["Value"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			placeholderString := controls["Placeholder"].Value.(string)

			return &InputText{
				Label: labelText,
				Value: valueString,
				Disabled: isDisabled,
				Placeholder: placeholderString,
			}
		},
	)

}

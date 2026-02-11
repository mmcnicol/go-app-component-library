//go:build dev
// pkg/components/select_one/select_one_stories.go
package select_one

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    /*
    selectOneDefault := &SelectOne{
        Label: "Label text.",
        Options: []string{"Go", "Python", "Rust", "JavaScript"},
    }

    storybook.Register("Select One", "Default", func() app.UI {
        return selectOneDefault
    })

    selectOneReadOnly := &SelectOne{
        Label: "Label text.",
        Options: []string{"Go", "Python", "Rust", "JavaScript"},
        Disabled: true,
    }

    storybook.Register("Select One", "ReadOnly", func() app.UI {
        return selectOneReadOnly
    })
    */

    selectOptions := []string{"Go", "Python", "Rust", "JavaScript"}

    storybook.Register("Form", "Select One", 
		map[string]*storybook.Control{
			"PromptText": {Label: "Prompt Text", Type: storybook.ControlText, Value: "Choose an option..."},
			"Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
			"Options": {Label: "Options", Type: storybook.ControlSelect, Options: selectOptions}, 
            "SelectedValue": {Label: "Selected Value", Type: storybook.ControlText, Value: ""}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			promptText := controls["PromptText"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
			opts := controls["Options"].Value.(string)
            selectedValue := controls["SelectedValue"].Value.(string)

			return &SelectOne{
				PromptText: promptText,
				Disabled: isDisabled,
				Options: opts,
                SelectedValue: selectedValue,
                OnSelect: func(ctx app.Context, val string) {
                    // Update the shared registry state
                    controls["SelectedValue"].Value = val
                    
                    // Refresh the UI so the Shell's Controls Panel sees the change
                    ctx.Update()
                },
			}
		},
	)
    
}

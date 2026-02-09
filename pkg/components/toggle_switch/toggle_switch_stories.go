//go:build dev
// pkg/components/toggle_switch/toggle_switch_stories.go
package toggle_switch

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    /*
    toggleSwitchDefault := &ToggleSwitch{
        IsOn:  false,
        Label: "Label text.",
    }

    storybook.Register("Toggle Switch", "Default", func() app.UI {
        return toggleSwitchDefault
    })

    toggleSwitchReadOnly := &ToggleSwitch{
        IsOn:  false,
        Label: "Label text.",
        Disabled: true,
    }

    storybook.Register("Toggle Switch", "ReadOnly", func() app.UI {
        return toggleSwitchReadOnly
    })
    */

    storybook.Register("Form", "Toggle Switch", 
		map[string]*storybook.Control{
			"On": {Label: "isOn", Type: storybook.ControlBool, Value: false},
			"Label": {Label: "Label", Type: storybook.ControlText, Value: "Label Text."}, 
            "Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
		},
		func(controls map[string]*storybook.Control) app.UI {
			isOn := controls["On"].Value.(bool)
            labelString := controls["Label"].Value.(string)
			isDisabled := controls["Disabled"].Value.(bool)
            
			return &ToggleSwitch{
                IsOn:  isOn,
                Label: labelString,
                Disabled: isDisabled,
                OnClick: func(ctx app.Context, val bool) {
                    // Update the shared registry state
                    controls["On"].Value = val
                    
                    // Refresh the UI so the Shell's Controls Panel sees the change
                    ctx.Update()
                },
            }
		},
	)
    
}

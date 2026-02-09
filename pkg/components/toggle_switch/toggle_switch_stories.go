//go:build dev
// pkg/components/toggle_switch/toggle_switch_stories.go
package toggle_switch

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

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
    
}

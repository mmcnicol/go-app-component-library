//go:build dev
// pkg/components/toggle_switch/toggle_switch_stories.go
package toggle_switch

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    toggleOff := &ToggleSwitch{
        IsOn:  false,
        Label: "Label text.",
    }

    storybook.Register("Toggle Switch", "Off", func() app.UI {
        return toggleOff
    })
    
    toggleOn := &ToggleSwitch{
        IsOn:  true,
        Label: "Label text.",
    }

    storybook.Register("Toggle Switch", "On", func() app.UI {
        return toggleOn
    })
}

//go:build dev
// pkg/components/button/button_stories.go
package button

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    storybook.Register("Form", "Button", 
        map[string]*storybook.Control{
            "Label":    {Label: "Label", Type: storybook.ControlText, Value: "Click Me"},
            "Disabled": {Label: "Disabled", Type: storybook.ControlBool, Value: false},
            "Look":     {
                Label: "Look", 
                Type: storybook.ControlSelect, 
                Options: []string{"primary", "secondary", "danger"},
            },
        },
        func(controls map[string]*storybook.Control) app.UI {
            return &Button{
                Label:    controls["Label"].Value.(string),
                Disabled: controls["Disabled"].Value.(bool),
                Look:     ButtonLook(controls["Look"].Options.(string)),
                OnClick: func(ctx app.Context, e app.Event) {
                    app.Log("Button Clicked!")
                },
            }
        },
    )

}

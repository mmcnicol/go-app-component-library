//go:build dev
// pkg/components/static_message/static_message_stories.go
package static_message

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

	storybook.Register("Messages", "StaticMessage", 
		map[string]*storybook.Control{
			"Severity": {Label: "Severity", Type: storybook.ControlText, Value: "info"}, 
			"Summary":  {Label: "Summary", Type: storybook.ControlText, Value: "Info Message"}, 
			"Detail":   {Label: "Detail", Type: storybook.ControlText, Value: "This is a detailed description of the alert."}, 
		},
		func(controls map[string]*storybook.Control) app.UI {
			severity := controls["Severity"].Value.(string)
			summary := controls["Summary"].Value.(string)
			detail := controls["Detail"].Value.(string)

			return &StaticMessage{
				Severity: severity,
				Summary:  summary,
				Detail:   detail,
			}
		},
	)

}


//go:build dev
// pkg/components/static_message/static_message_stories.go
package static_message

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {

	storybook.Register("Messages", "StaticMessage", 
		map[string]*storybook.Control{
			"Severity": {
				Label: "Severity", 
				Type: storybook.ControlSelect, 
				Value: "info",
				Options: []string{"info", "success", "warn", "error"},
			}, 
			"Summary": {
				Label: "Summary", 
				Type: storybook.ControlText, 
				Value: "Information Message",
			}, 
			"Detail": {
				Label: "Detail", 
				Type: storybook.ControlText, 
				Value: "This is a detailed description of the alert.",
			}, 
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

	/*
	// Optional: Add individual stories for each severity type as examples
	storybook.Register("Messages", "Info Message", 
		map[string]*storybook.Control{},
		func(controls map[string]*storybook.Control) app.UI {
			return &StaticMessage{
				Severity: "info",
				Summary:  "Information",
				Detail:   "Your request has been processed successfully.",
			}
		},
	)
	
	storybook.Register("Messages", "Success Message", 
		map[string]*storybook.Control{},
		func(controls map[string]*storybook.Control) app.UI {
			return &StaticMessage{
				Severity: "success",
				Summary:  "Success!",
				Detail:   "Your changes have been saved successfully.",
			}
		},
	)
	
	storybook.Register("Messages", "Warning Message", 
		map[string]*storybook.Control{},
		func(controls map[string]*storybook.Control) app.UI {
			return &StaticMessage{
				Severity: "warn",
				Summary:  "Warning",
				Detail:   "Please review your information before submitting.",
			}
		},
	)
	
	storybook.Register("Messages", "Error Message", 
		map[string]*storybook.Control{},
		func(controls map[string]*storybook.Control) app.UI {
			return &StaticMessage{
				Severity: "error",
				Summary:  "Error",
				Detail:   "An error occurred while processing your request.",
			}
		},
	)
	*/

}

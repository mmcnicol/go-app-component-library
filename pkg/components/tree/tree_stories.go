//go:build dev
// pkg/components/tree/tree_stories.go
package tree

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {
	if app.IsClient {
		app.Log("Tree init()")
	}

	storybook.Register("Data", "Tree", 
		map[string]*storybook.Control{
			"Expanded": {Label: "Expand All", Type: storybook.ControlBool, Value: false},
		},
		func(controls map[string]*storybook.Control) app.UI {
			expandAll := controls["Expanded"].Value.(bool)

			// Sample clinical data structure
			data := []*TreeNode{
				{
					Label: "Patient Records",
					Icon:  "folder",
					Expanded: expandAll,
					Children: []*TreeNode{
						{
							Label: "Lab Results",
							Icon:  "folder",
							Children: []*TreeNode{
								{Label: "Biochemistry_2024.pdf", Icon: "file"},
								{Label: "Hematology_FullCount.xml", Icon: "file"},
							},
						},
						{
							Label: "Discharge Summaries",
							Icon:  "folder",
							Children: []*TreeNode{
								{Label: "Final_Discharge_JohnDoe.pdf", Icon: "file"},
							},
						},
					},
				},
			}

			return app.Div().Style("padding", "20px").Body(
				&Tree{Data: data},
			)
		},
	)

}

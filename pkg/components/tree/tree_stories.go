//go:build dev
// pkg/components/tree/tree_stories.go
package tree

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Create a persistent variable for the story data
var treeData = []*TreeNode{
    {
        Label: "Patient Records",
        Icon:  "folder",
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

func init() {
    storybook.Register("Data", "Tree", 
        map[string]*storybook.Control{
            "Selected": {Label: "Active Document", Type: storybook.ControlText, Value: "None", ReadOnly: true},
        },
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Style("padding", "20px").Body(
                &Tree{
                    Data: treeData,
                    // Pass the controls so the component can update the sidebar
                    OnSelect: func(ctx app.Context, nodeName string) {
                        controls["Selected"].Value = nodeName
                        ctx.Update()
                    },
                },
            )
        },
    )
}

/*
func init() {
    storybook.Register("Data", "Tree", 
        map[string]*storybook.Control{
            //"Expanded": {Label: "Expand All", Type: storybook.ControlBool, Value: false},
        },
        func(controls map[string]*storybook.Control) app.UI {
            //expandAll := controls["Expanded"].Value.(bool)

            // If "Expand All" is toggled in the sidebar, update the persistent data
            //if expandAll {
            //    setAllExpanded(treeData, true)
            //}

            return app.Div().Style("padding", "20px").Body(
                &Tree{Data: treeData},
            )
        },
    )
}
*/

/*
// Helper to handle the "Expand All" control logic
func setAllExpanded(nodes []*TreeNode, state bool) {
    for _, n := range nodes {
        n.Expanded = state
        if len(n.Children) > 0 {
            setAllExpanded(n.Children, state)
        }
    }
}
*/

/*
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
*/

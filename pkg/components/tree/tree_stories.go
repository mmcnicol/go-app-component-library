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
            "Selected": {
				Label: "Active Document", 
				Type: storybook.ControlText, 
				Value: "None", ReadOnly: true,
			},
        },
        func(controls map[string]*storybook.Control) app.UI {
            return app.Div().Style("padding", "20px").Body(
                &Tree{
                    Data: treeData,
                    // Pass the controls so the component can update the sidebar
                    OnSelect: func(ctx app.Context, nodeName string) {
                        controls["Selected"].Value = nodeName
						app.Log("Document selected: " + nodeName)
                        ctx.Update()
                    },
                },
            )
        },
    )
}

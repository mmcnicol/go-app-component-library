// pkg/components/tree/tree.go
package tree

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/components/icon"
	"github.com/mmcnicol/go-app-component-library/pkg/theme"
)

type TreeNode struct {
	Label    string
	Expanded bool
	Selected bool
	Children []*TreeNode
	Icon     string
}

type Tree struct {
	app.Compo
	Data []*TreeNode
	OnSelect func(ctx app.Context, nodeName string) // Callback for parent sync
}

func (t *Tree) Render() app.UI {
	if app.IsClient {
		app.Log("Tree Render()")
	}
	return app.Div().Class("ui-tree").Body(
		app.Range(t.Data).Slice(func(i int) app.UI {
			return t.renderNode(t.Data[i], 0)
		}),
	)
}

func (t *Tree) renderNode(node *TreeNode, level int) app.UI {
	i := &icon.Icon{}
	hasChildren := len(node.Children) > 0

	// Determine class based on selection state
	nodeClass := "tree-component-node"
	if node.Selected {
		nodeClass = "tree-component-node tree-component-node-active"
	}

	return app.Div().Body(
		app.Div().
			Class(nodeClass). // Apply the CSS class here
			Style("display", "flex").
			Style("align-items", "center").
			//Style("padding", "6px 8px").
			Style("padding", theme.SpacingSM + " 0"). // Use spacing tokens
			Style("margin", "2px 0").
			//Style("padding-left", app.FormatString("%dpx", level*32)).
			Style("padding-left", app.FormatString("calc(%d * %s)", level, theme.SpacingXL)).
			Style("cursor", "pointer").
			Style("border-radius", "4px").
			OnClick(func(ctx app.Context, e app.Event) {
				if hasChildren {
					node.Expanded = !node.Expanded
				} else {
					t.deselectAll(t.Data)
					node.Selected = true
					if t.OnSelect != nil {
						t.OnSelect(ctx, node.Label)
					}
				}
				ctx.Update()
			}).
			Body(
				app.If(hasChildren, func() app.UI {
					if node.Expanded {
						return i.GetIcon("chevron-down", 16)
					}
					return i.GetIcon("chevron-right", 16)
				}).Else(func() app.UI {
					return app.Div().Style("width", "16px")
				}),
				
				// Icon
				app.If(node.Icon != "", func() app.UI {
					return app.Div().Style("margin", "0 6px").Body(i.GetIcon(node.Icon, 18))
				}),
				
				//app.Span().Style("font-weight", "500").Text(node.Label),
				// Label with Weight - CSS handles the truncation now
                app.Span().
                    Style("font-weight", "500").
                    Text(node.Label).
                    Title(node.Label), // Browser tooltip shows the full name on hover
			),
		
		app.If(node.Expanded && hasChildren, func() app.UI {
			return app.Div().Body(
				app.Range(node.Children).Slice(func(idx int) app.UI {
					return t.renderNode(node.Children[idx], level+1)
				}),
			)
		}),
	)
}

/*
func (t *Tree) renderNode(node *TreeNode, level int) app.UI {
    i := &icon.Icon{}
    hasChildren := len(node.Children) > 0

    // Define the selection color (Clinical Blue)
    bg := "transparent"
    textColor := "inherit"
    if node.Selected {
        bg = "#E3F2FD"    // Light blue background
        textColor = "#0D47A1" // Dark blue text
    }

    return app.Div().Body(
        app.Div().
            Style("display", "flex").
            Style("align-items", "center").
            Style("padding", "6px 8px").
            Style("margin", "2px 0").
            Style("padding-left", app.FormatString("%dpx", level*32)).
            Style("cursor", "pointer").
            Style("background-color", bg).
            Style("color", textColor).
            Style("border-radius", "4px").
			OnClick(func(ctx app.Context, e app.Event) {
                if hasChildren {
                    node.Expanded = !node.Expanded
                } else {
                    t.deselectAll(t.Data)
                    node.Selected = true
					//app.Log("Document selected: " + node.Label)
                    
                    // Trigger the callback to update Storybook controls
                    if t.OnSelect != nil {
                        t.OnSelect(ctx, node.Label)
                    }
                }
                ctx.Update()
            }).
            Body(
                // Toggle Icon (Chevron)
                app.If(hasChildren, func() app.UI {
                    if node.Expanded {
                        return i.GetIcon("chevron-down", 16)
                    }
                    return i.GetIcon("chevron-right", 16)
                }).Else(func() app.UI {
                    return app.Div().Style("width", "16px")
                }),
                
                // File/Folder Icon
                app.If(node.Icon != "", func() app.UI {
                    return app.Div().Style("margin", "0 6px").Body(i.GetIcon(node.Icon, 18))
                }),
                
                app.Span().Style("font-weight", "500").Text(node.Label),
            ),
        
        // Recursive Children
        app.If(node.Expanded && hasChildren, func() app.UI {
            return app.Div().Body(
                app.Range(node.Children).Slice(func(idx int) app.UI {
                    return t.renderNode(node.Children[idx], level+1)
                }),
            )
        }),
    )
}
*/

// Helper to clear existing selections
func (t *Tree) deselectAll(nodes []*TreeNode) {
    for _, n := range nodes {
        n.Selected = false
        if len(n.Children) > 0 {
            t.deselectAll(n.Children)
        }
    }
}

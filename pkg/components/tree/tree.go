// pkg/components/tree/tree.go
package tree

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/components/icon"
)

type TreeNode struct {
	Label    string
	Expanded bool
	Children []*TreeNode
	Icon     string
}

type Tree struct {
	app.Compo
	Data []*TreeNode
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

	return app.Div().Body(
		app.Div().
			Style("display", "flex").
			Style("align-items", "center").
			Style("padding", "4px 0").
			Style("padding-left", app.FormatString("%dpx", level*20)).
			Style("cursor", "pointer").
			OnClick(func(ctx app.Context, e app.Event) {
				if app.IsClient {
					app.Log("Tree OnClick()")
				}
				node.Expanded = !node.Expanded
				ctx.Update()
			}).
			Body(
				// Toggle Icon
				app.If(hasChildren, func() app.UI {
					if node.Expanded {
						return i.GetIcon("chevron-down", 16)
					}
					return i.GetIcon("chevron-right", 16)
				}).Else(func() app.UI {
					return app.Div().Style("width", "16px") // Spacer
				}),
				
				// Node Icon (optional)
				app.If(node.Icon != "", func() app.UI {
					return app.Div().Style("margin", "0 4px").Body(i.GetIcon(node.Icon, 18))
				}),
				
				app.Span().Text(node.Label),
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

//go:build dev

package chart

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
	storybook.Register("Charts", "Streaming Chart",
		map[string]*storybook.Control{
			"Buffer Capacity": storybook.NewRangeControl(50, 500, 50, 100),
			"Line Color":      storybook.NewColorControl("#10b981"),
		},
		func(controls map[string]*storybook.Control) app.UI {
			cap := controls["Buffer Capacity"].Value.(int)
			
			return New(nil, // No initial data needed for streaming
				WithTitle("Live Data Feed"),
				WithColor(controls["Line Color"].Value.(string)),
				SetStreaming(cap),
			)
		},
	)
}

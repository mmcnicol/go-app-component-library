//go:build dev
// pkg/components/select_one/select_one_stories.go
package select_one

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

// Use init() to auto-register when this package is imported
func init() {

    selectOne := &SelectOne{
        Label: "Label text.",
    }

    storybook.Register("Select One", "Default", func() app.UI {
        return selectOne
    })
    
}

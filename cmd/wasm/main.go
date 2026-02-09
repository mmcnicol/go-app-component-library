//go:build dev

package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
	
    // Import your components so their init() functions run and register stories
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/hello"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/phase-banner"
)

func main() {
	// Route "/" to the Storybook shell
	app.Route("/", func() app.Composer { return &storybook.Shell{} })

	app.RunWhenOnBrowser()
}

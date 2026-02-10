//go:build dev
// cmd/wasm/main.go
package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/storybook"
	
	// Import your components so their init() functions run and register stories
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/built_in"
	//_ "github.com/mmcnicol/go-app-component-library/pkg/components/hello"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/phase_banner"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/toggle_switch"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/select_one"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/input_text"
	_ "github.com/mmcnicol/go-app-component-library/pkg/components/static_message"
)

func main() {
	// Route "/" to the Storybook shell
	app.Route("/", func() app.Composer { return &storybook.Shell{} })

	app.RunWhenOnBrowser()
}

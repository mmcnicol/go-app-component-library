// cmd/wasm/main.go
package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"../../pkg/components"
)

func main() {
	
	// Register the components that correspond to routes
	app.Route("/", func() app.Composer { return &components.Hello{} })

	// This function starts the Wasm app in the browser.
	// It stays idle when running on the server.
	app.RunWhenOnBrowser()
}

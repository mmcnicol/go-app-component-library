package main

import (
	"encoding/json"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/mmcnicol/go-app-component-library/pkg/components"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// CRITICAL: The server MUST know the route exists
	// Register the components that correspond to routes
	app.Route("/", func() app.Composer { return &components.Hello{} })

	h := &app.Handler{
		Name:      "go-app component library",
		Description: "A go-app UI library using Go and WebAssembly",
		Author:      "mmcnicol",
		Styles: []string{
			//"/web/css/main.css", // Link to your custom CSS
		},
		Icon: app.Icon{
			//Default: "/web/images/logo.png",
		},
		//Resources: app.LocalDir("web"),
	}

	http.Handle("/", h)

	// Example API endpoint
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
	})

	log.Println("Serving at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

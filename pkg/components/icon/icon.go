// pkg/components/icon/icon.go
package icon

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Icon struct {
	app.Compo
}

// GetIcon returns an SVG icon
func (i *Icon) GetIcon(name string, size int) app.UI {
	var d string
	isSpinner := name == "spinner"

	switch name {
	case "success":
		d = "M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"
	case "info":
		d = "M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"
	case "warn":
		d = "M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"
	case "error":
		d = "M12 2C6.47 2 2 6.47 2 12s4.47 10 10 10 10-4.47 10-10S17.53 2 12 2zm5 13.59L15.59 17 12 13.41 8.41 17 7 15.59 10.59 12 7 8.41 8.41 7 12 10.59 15.59 7 17 8.41 13.41 12 17 15.59z"
	case "chevron-right":
		d = "M10 6L8.59 7.41 13.17 12l-4.58 4.59L10 18l6-6z"
	case "chevron-down":
		d = "M16.59 8.59L12 13.17 7.41 8.59 6 10l6 6 6-6z"
	case "chevron-left":
		d = "M15.41 7.41L14 6l-6 6 6 6 1.41-1.41L10.83 12z"
	case "chevron-up":
		d = "M7.41 15.41L12 10.83l4.59 4.58L18 14l-6-6-6 6z"
	case "sort":
		d = "M3 18h6v-2H3v2zM3 6v2h18V6H3zm0 7h12v-2H3v2z"
	case "settings":
		d = "M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.07-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.07.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"
	case "logout":
		d = "M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"
	case "folder":
		d = "M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"
	case "file":
		d = "M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"
	/*
	case "hospital-inpatient":
		// Consistent: Hospital bed (inpatient) with medical weight
		d = "M18 10h-5V5h-2v5H6v2h5v5h2v-5h5v-2zm3 7H3v-2h18v2zm-2-5h-2v3H7v-3H5v6h14v-6z"
	case "hospital-outpatient":
		// Consistent: Medical box/cross with a walking figure (ambulatory)
		d = "M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-4v4h-2v-4H8v-2h4V8h2v4h4v2zM7 7c1.1 0 2 .9 2 2s-.9 2-2 2-2-.9-2-2 .9-2 2-2zm-2 9v-1c0-1.1.9-2 2-2s2 .9 2 2v1h-4z"
	case "hospital-alert-covid":
		// Improved: Medical cross with a central 'virus' alert point to match the outpatient style
		d = "M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-1 11h-3.12c-.22-.61-.69-1.12-1.32-1.38V9h2v2h2v2h-2v1zm-6-5h2v2.62c-.63.26-1.1.77-1.32 1.38H9V11h2V9zm-2 4h3.12c.22.61.69 1.12 1.32 1.38V18h-2v-2H9v-2h2v-1zm6 5h-2v-2.62c.63-.26 1.1-.77 1.32-1.38H18v2h-2v2h-1z"
	*/
	case "hospital-inpatient":
		// A clean, recognizable hospital bed with a medical cross above it
		d = "M19 7h-8v7H3V5H1v15h2v-3h18v3h2v-9c0-2.21-1.79-4-4-4zm-7 5h2V9h2v3h2v2h-2v3h-2v-3h-2v-2z"
	case "hospital-outpatient":
		// A medical cross with a person silhouette next to it (Ambulatory)
		d = "M18 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm-11 2v-3h2v3h3v2H9v3H7v-3H4v-2h3zm11 1c-2.67 0-8 1.34-8 4v3h16v-3c0-2.66-5.33-4-8-4z"
	case "pharmacy":
		// A mortar and pestle—the universal symbol for pharmacy/medication
		d = "M21 5h-2.64l1.14-3.14L17.59 1 16.2 4.81 3 5v2h2l2 6-2 6v2h14v-2l-2-6 2-6h2V5zM7 17l1.33-4h7.34L17 17H7z"
	case "emergency":
		// A siren/beacon to represent urgent care or emergency services
		d = "M7 11h2V7H7v4zm10-4h-2v4h2V7zm-6 4h2V7h-2v4zM12 2C6.48 2 2 6.48 2 12v9h20v-9c0-5.52-4.48-10-10-10zm8 17H4v-7c0-4.41 3.59-8 8-8s8 3.59 8 8v7z"
	case "ambulance":
		// A medical vehicle silhouette with the distinctive cross
		d = "M19 11h-3V8h-2v3h-3v2h3v3h2v-3h3v-2zm3-3c0-1.1-.9-2-2-2h-3V4h-2v2H5c-1.1 0-2 .9-2 2v8c0 1.1.9 2 2 2h1.1c.45 1.18 1.58 2 2.9 2s2.45-.82 2.9-2h4.2c.45 1.18 1.58 2 2.9 2s2.45-.82 2.9-2H21c1.1 0 2-.9 2-2V8zm-13 9c-.55 0-1-.45-1-1s.45-1 1-1 1 .45 1 1-.45 1-1 1zm9 0c-.55 0-1-.45-1-1s.45-1 1-1 1 .45 1 1-.45 1-1 1z"		
	case "stethoscope":
		// Consultation icon: Stethoscope with a clean, circular chest piece
		d = "M12 2a5 5 0 00-5 5v5a5 5 0 0010 0V7a5 5 0 00-5-5zm3 10a3 3 0 01-6 0V7a3 3 0 016 0v5zm-3 8a3 3 0 100-6 3 3 0 000 6zm6-6v-1h2v1c0 3.87-3.13 7-7 7s-7-3.13-7-7v-1h2v1c0 2.76 2.24 5 5 5s5-2.24 5-5z"
	case "records":
		// Medical records icon: Clipboard with a centered medical cross
		d = "M18 2h-3.18C14.4 0.84 13.3 0 12 0s-2.4.84-2.82 2H6c-1.1 0-2 .9-2 2v16c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-6 0c.55 0 1 .45 1 1s-.45 1-1 1-1-.45-1-1 .45-1 1-1zm3 14h-2v2h-2v-2H9v-2h2v-2h2v2h2v2z"
	case "telehealth":
		// A computer monitor featuring a centered medical cross
		d = "M21 2H3c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h7l-2 3v1h8v-1l-2-3h7c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm0 14H3V4h18v12zm-5-7h-3V6h-2v3H8v2h3v3h2v-3h3V9z"
	case "vaccine":
		// A syringe icon with a visible droplet and plunger
		d = "M19.88 3.5c-.39-.39-1.03-.39-1.42 0l-1.06 1.06-1.54-1.54-1.41 1.41 1.54 1.54-1.77 1.77-2.12-2.12-1.41 1.41 2.12 2.12-4.24 4.24H6v2l2.3 2.3L3 22.1l1.41 1.41 4.51-4.51 2.3 2.3h2v-2.41l4.24-4.24 2.12 2.12 1.41-1.41-2.12-2.12 1.77-1.77 1.54 1.54 1.41-1.41-1.54-1.54 1.06-1.06c.39-.39.39-1.03 0-1.42z"
	case "hospital-alert-covid":
		// A medical cross with a high-visibility alert '!' in the center
		d = "M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-6 15h-2v-2h2v2zm0-4h-2V7h2v7z"
	case "spinner":
		d = "M12 4V2C6.48 2 2 6.48 2 12h2c0-4.41 3.59-8 8-8z"
	}

	classNames := fmt.Sprintf("icon icon-%s", name)
	if isSpinner {
		classNames += " icon-spin"
	}

	// Using app.Raw avoids the "undefined: app.Svg" compilation error
	// ViewBox stays 0 0 24 24. Width/Height handles the zoom.
	return app.Raw(fmt.Sprintf(
		`<svg fill="currentColor" class="%s" viewBox="0 0 24 24" width="%d" height="%d"><path d="%s" /></svg>`,
		classNames, size, size, d,
	))
}

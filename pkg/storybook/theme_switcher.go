// pkg/storybook/theme_switcher.go
package storybook

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ThemeSwitcher struct {
	app.Compo
	IsDark   bool
	OnChange func(ctx app.Context, isDark bool)
}

func (s *ThemeSwitcher) Render() app.UI {
	icon := "☀️"
	if s.IsDark {
		icon = "🌙"
	}

	return app.Div().
		Style("cursor", "pointer").
		OnClick(func(ctx app.Context, e app.Event) {
			s.IsDark = !s.IsDark
			if s.OnChange != nil {
				s.OnChange(ctx, s.IsDark) 
			}
		}).
		Body(
			app.Span().Text(icon),
		)
}

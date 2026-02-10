// pkg/components/input_text_area/input_text_area.go
package input_text_area

type InputTextArea struct {
	app.Compo
	Value       string
	Placeholder string
	Rows        int
	Cols        int
	Disabled    bool
	ReadOnly    bool
	OnInput     func(app.Context, app.Event)
}

func (t *InputTextArea) Render() app.UI {
	return app.Textarea().
		Class("input-textarea").
		Value(t.Value).
		Placeholder(t.Placeholder).
		Rows(t.Rows).
		Cols(t.Cols).
		Disabled(t.Disabled).
		ReadOnly(t.ReadOnly).
		OnInput(t.OnInput).
		Style("width", "100%").
		Style("padding", "8px").
		Style("border-radius", "4px").
		Style("border", "1px solid #ccc")
}

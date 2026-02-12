
github.com/mmcnicol/go-app-component-library/pkg/components/vizzy/

a chart/visulization library

goals:
* Separation of Concerns: Chart logic vs rendering logic
* Testability: Pure data models are easy to test
* Extensibility: Easy to add new renderers (SVG, WebGL, etc.)
* Performance: Can optimize renderers independently
* Maintainability: Smaller, focused files
* Reusability: Same chart model can be rendered by different engines

the primary renderer is for HTML5 canvas as a go-app v10 component.

the package has a flat file structure.

file:
core_chart.go             # Pure data models
core_registry.go          # Engine registry
core_types.go             # Shared types
core_interfaces.go        # Core interfaces
render_canvas_engine.go   # Canvas engine implementation
render_canvas_bar.go      # Bar chart renderer
render_canvas_line.go     # Line chart renderer
render_canvas_...
render_svg_...
component_chart.go        # app.UI wrapper


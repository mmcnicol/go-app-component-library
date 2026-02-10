// pkg/components/table/table.go
package table

import (
    "fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// TableProps defines properties for the base table component
type TableProps struct {
    ID           string
    Class        string
    Style        map[string]string
    Caption      string
    Columns      []Column
    Data         []map[string]interface{}
    RowKey       string // Field to use as unique row identifier
    Striped      bool
    Bordered     bool
    Hoverable    bool
    Compact      bool
    Responsive   bool
    OnRowClick   func(ctx app.Context, rowData map[string]interface{}, index int)
    OnRowHover   func(ctx app.Context, rowData map[string]interface{}, index int)
    EmptyState   app.UI
    Loading      bool
    LoadingState app.UI
    DataTestID   string
}

// Column defines a table column configuration
type Column struct {
    ID           string
    Header       string
    Accessor     string // Key to access data in row object
    Width        string
    MinWidth     string
    MaxWidth     string
    Align        string // "left", "center", "right"
    Sortable     bool
    Filterable   bool
    CellRenderer func(data interface{}, rowIndex int, colIndex int) app.UI
    HeaderRenderer func(col Column, colIndex int) app.UI
    FooterRenderer func(col Column, colIndex int) app.UI
}

// Table represents the base table component
type Table struct {
    app.Compo
    props TableProps
}

// Render returns the UI representation of the table
func (t *Table) Render() app.UI {
    // Build table classes
    class := "table"
    
    if t.props.Striped {
        class += " table--striped"
    }
    
    if t.props.Bordered {
        class += " table--bordered"
    }
    
    if t.props.Hoverable {
        class += " table--hoverable"
    }
    
    if t.props.Compact {
        class += " table--compact"
    }
    
    if t.props.Class != "" {
        class += " " + t.props.Class
    }
    
    // Create table element
    /*
    table := app.Table().
        ID(t.props.ID).
        Class(class).
        Style(t.props.Style).
        DataSet("testid", t.props.DataTestID)
    */
    table := app.Table().
        ID(t.props.ID).
        Class(class)
    
    // Apply styles if provided
    if len(t.props.Style) > 0 {
        var styleArgs []string
        for k, v := range t.props.Style {
            styleArgs = append(styleArgs, k, v)
        }
        table = table.Style(styleArgs...)
    }

    table = table.DataSet("testid", t.props.DataTestID)

    // Add caption if provided
    if t.props.Caption != "" {
        table = table.Body(
            app.Caption().Class("table__caption").Text(t.props.Caption),
        )
    }
    
    // Add responsive wrapper if needed
    var content app.UI
    if t.props.Responsive {
        content = app.Div().
            Class("table-responsive").
            Body(
                table.Body(
                    t.renderHeader(),
                    t.renderBody(),
                    t.renderFooter(),
                ),
            )
    } else {
        content = table.Body(
            t.renderHeader(),
            t.renderBody(),
            t.renderFooter(),
        )
    }
    
    // Handle loading and empty states
    if t.props.Loading {
        if t.props.LoadingState != nil {
            return t.props.LoadingState
        }
        return t.renderLoadingState()
    }
    
    if len(t.props.Data) == 0 {
        if t.props.EmptyState != nil {
            return t.props.EmptyState
        }
        return t.renderEmptyState()
    }
    
    return content
}

func (t *Table) renderHeader() app.UI {
    var headers []app.UI
    
    for i, col := range t.props.Columns {
        headerClass := "table__header"
        
        if col.Align != "" {
            headerClass += " table__header--align-" + col.Align
        }
        
        var headerContent app.UI
        if col.HeaderRenderer != nil {
            headerContent = col.HeaderRenderer(col, i)
        } else {
            headerContent = app.Text(col.Header)
        }
        
        /*
        headers = append(headers, app.Th().
            ID(col.ID).
            Class(headerClass).
            Style(t.getColumnStyle(col)...). // Use spread operator
            Scope("col").
            Body(headerContent))
        */

        styleArgs := t.getColumnStyle(col)
        if len(styleArgs) > 0 {
            headers = append(headers, app.Th().
                ID(col.ID).
                Class(headerClass).
                Style(styleArgs...). // Use spread operator
                Scope("col").
                Body(headerContent))
        } else {
            headers = append(headers, app.Th().
                ID(col.ID).
                Class(headerClass).
                Scope("col").
                Body(headerContent))
        }
    }
    
    return app.THead().
        Class("table__head").
        Body(
            app.Tr().Class("table__row").Body(headers...),
        )
}

func (t *Table) renderBody() app.UI {
    var rows []app.UI
    
    for rowIdx, rowData := range t.props.Data {
        rowKey := t.getRowKey(rowData, rowIdx)
        rowClass := "table__row"
        
        if t.props.Hoverable {
            rowClass += " table__row--hoverable"
        }
        
        if t.props.OnRowClick != nil {
            rowClass += " table__row--clickable"
        }
        
        row := app.Tr().
            ID(rowKey).
            Class(rowClass).
            DataSet("row-index", rowIdx).
            OnClick(t.handleRowClick(rowData, rowIdx)).
            OnMouseOver(t.handleRowHover(rowData, rowIdx))
        
        var cells []app.UI
        
        for colIdx, col := range t.props.Columns {
            cellClass := "table__cell"
            
            if col.Align != "" {
                cellClass += " table__cell--align-" + col.Align
            }
            
            var cellContent app.UI
            if col.CellRenderer != nil {
                cellData := rowData[col.Accessor]
                cellContent = col.CellRenderer(cellData, rowIdx, colIdx)
            } else {
                cellData := rowData[col.Accessor]
                cellContent = t.formatCellValue(cellData)
            }
            
            /*
            cells = append(cells, app.Td().
                Class(cellClass).
                Style(t.getColumnStyle(col)...). // Use spread operator
                Body(cellContent))
            */

            styleArgs := t.getColumnStyle(col)
            if len(styleArgs) > 0 {
                cells = append(cells, app.Td().
                    Class(cellClass).
                    Style(styleArgs...). // Use spread operator
                    Body(cellContent))
            } else {
                cells = append(cells, app.Td().
                    Class(cellClass).
                    Body(cellContent))
            }
        }
        
        rows = append(rows, row.Body(cells...))
    }
    
    return app.TBody().
        Class("table__body").
        Body(rows...)
}

func (t *Table) renderFooter() app.UI {
    // Only render footer if at least one column has a footer renderer
    hasFooter := false
    for _, col := range t.props.Columns {
        if col.FooterRenderer != nil {
            hasFooter = true
            break
        }
    }
    
    if !hasFooter {
        return app.TFoot()
    }
    
    var footers []app.UI
    
    for i, col := range t.props.Columns {
        footerClass := "table__footer"
        
        if col.Align != "" {
            footerClass += " table__footer--align-" + col.Align
        }
        
        var footerContent app.UI
        if col.FooterRenderer != nil {
            footerContent = col.FooterRenderer(col, i)
        }
        
        /*
        footers = append(footers, app.Td().
            Class(footerClass).
            Style(t.getColumnStyle(col)...). // Use spread operator
            Body(footerContent))
        */

        styleArgs := t.getColumnStyle(col)
        if len(styleArgs) > 0 {
            footers = append(footers, app.Td().
                Class(footerClass).
                Style(styleArgs...). // Use spread operator
                Body(footerContent))
        } else {
            footers = append(footers, app.Td().
                Class(footerClass).
                Body(footerContent))
        }
    }
    
    return app.TFoot().
        Class("table__foot").
        Body(
            app.Tr().Class("table__row").Body(footers...),
        )
}

func (t *Table) getRowKey(rowData map[string]interface{}, index int) string {
    if t.props.RowKey != "" {
        if key, ok := rowData[t.props.RowKey].(string); ok {
            return key
        }
    }
    return "row-" + string(index)
}

// Update getColumnStyle to properly handle Style() method
func (t *Table) getColumnStyle(col Column) []any {
    var styles []any
    
    if col.Width != "" {
        styles = append(styles, "width", col.Width)
    }
    
    if col.MinWidth != "" {
        styles = append(styles, "min-width", col.MinWidth)
    }
    
    if col.MaxWidth != "" {
        styles = append(styles, "max-width", col.MaxWidth)
    }
    
    return styles
}

func (t *Table) formatCellValue(value interface{}) app.UI {
    switch v := value.(type) {
    case string:
        return app.Text(v)
    case int, int32, int64, float32, float64:
        return app.Text(fmt.Sprintf("%v", v))
    case bool:
        if v {
            return app.Raw(`<span class="icon-check"></span>`)
        }
        return app.Raw(`<span class="icon-x"></span>`)
    case nil:
        return app.Text("-")
    default:
        return app.Text(fmt.Sprintf("%v", v))
    }
}

func (t *Table) handleRowClick(rowData map[string]interface{}, index int) func(ctx app.Context, e app.Event) {
    return func(ctx app.Context, e app.Event) {
        if t.props.OnRowClick != nil {
            t.props.OnRowClick(ctx, rowData, index)
        }
    }
}

func (t *Table) handleRowHover(rowData map[string]interface{}, index int) func(ctx app.Context, e app.Event) {
    return func(ctx app.Context, e app.Event) {
        if t.props.OnRowHover != nil {
            t.props.OnRowHover(ctx, rowData, index)
        }
    }
}

func (t *Table) renderLoadingState() app.UI {
    return app.Div().
        Class("table__loading").
        Body(
            app.Div().Class("spinner"),
            app.P().Text("Loading data..."),
        )
}

func (t *Table) renderEmptyState() app.UI {
    return app.Div().
        Class("table__empty").
        Body(
            app.Div().Class("table__empty-icon"),
            app.H3().Text("No data available"),
            app.P().Text("There is no data to display at the moment."),
        )
}

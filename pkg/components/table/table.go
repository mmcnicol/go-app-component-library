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
    // Create table element with styles
    table := t.createTableElement()

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
        
        headers = append(headers, t.createHeaderCell(col, headerClass, headerContent))
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
            
            cells = append(cells, t.createBodyCell(col, cellClass, cellContent))
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
        
        footers = append(footers, t.createFooterCell(col, footerClass, footerContent))
    }
    
    return app.TFoot().
        Class("table__foot").
        Body(
            app.Tr().Class("table__row").Body(footers...),
        )
}

// Helper function to apply table styles
func (t *Table) applyTableStyles(elem app.HTMLTag) app.HTMLTag {
    for k, v := range t.props.Style {
        elem = elem.Style(k, v)
    }
    return elem
}

// Helper function to apply column styles
func (t *Table) applyColumnStyles(elem app.HTMLTag, col Column) app.HTMLTag {
    if col.Width != "" {
        elem = elem.Style("width", col.Width)
    }
    if col.MinWidth != "" {
        elem = elem.Style("min-width", col.MinWidth)
    }
    if col.MaxWidth != "" {
        elem = elem.Style("max-width", col.MaxWidth)
    }
    return elem
}

/*
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
        
        headerElem := app.Th().
            ID(col.ID).
            Class(headerClass).
            Scope("col")
        
        // Apply column styles
        headerElem = t.applyColumnStyles(headerElem, col)
        
        headers = append(headers, headerElem.Body(headerContent))
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
            
            cellElem := app.Td().Class(cellClass)
            
            // Apply column styles
            cellElem = t.applyColumnStyles(cellElem, col)
            
            cells = append(cells, cellElem.Body(cellContent))
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
        
        footerElem := app.Td().Class(footerClass)
        
        // Apply column styles
        footerElem = t.applyColumnStyles(footerElem, col)
        
        footers = append(footers, footerElem.Body(footerContent))
    }
    
    return app.TFoot().
        Class("table__foot").
        Body(
            app.Tr().Class("table__row").Body(footers...),
        )
}
*/

func (t *Table) getRowKey(rowData map[string]interface{}, index int) string {
    if t.props.RowKey != "" {
        if key, ok := rowData[t.props.RowKey].(string); ok {
            return key
        }
    }
    return "row-" + string(index)
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

// Helper function to apply table styles during creation
func (t *Table) createTableElement() app.UI {
    table := app.Table().
        ID(t.props.ID).
        Class(t.getTableClasses())
    
    // Apply styles if provided
    if len(t.props.Style) > 0 {
        for k, v := range t.props.Style {
            table = table.Style(k, v)
        }
    }
    
    return table.DataSet("testid", t.props.DataTestID)
}

// Helper function to get table classes
func (t *Table) getTableClasses() string {
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
    
    return class
}

// Helper function to create header cell with styles
func (t *Table) createHeaderCell(col Column, headerClass string, headerContent app.UI) app.UI {
    headerElem := app.Th().
        ID(col.ID).
        Class(headerClass).
        Scope("col")
    
    // Apply column styles
    if col.Width != "" {
        headerElem = headerElem.Style("width", col.Width)
    }
    if col.MinWidth != "" {
        headerElem = headerElem.Style("min-width", col.MinWidth)
    }
    if col.MaxWidth != "" {
        headerElem = headerElem.Style("max-width", col.MaxWidth)
    }
    
    return headerElem.Body(headerContent)
}

// Helper function to create body cell with styles
func (t *Table) createBodyCell(col Column, cellClass string, cellContent app.UI) app.UI {
    cellElem := app.Td().Class(cellClass)
    
    // Apply column styles
    if col.Width != "" {
        cellElem = cellElem.Style("width", col.Width)
    }
    if col.MinWidth != "" {
        cellElem = cellElem.Style("min-width", col.MinWidth)
    }
    if col.MaxWidth != "" {
        cellElem = cellElem.Style("max-width", col.MaxWidth)
    }
    
    return cellElem.Body(cellContent)
}

// Helper function to create footer cell with styles
func (t *Table) createFooterCell(col Column, footerClass string, footerContent app.UI) app.UI {
    footerElem := app.Td().Class(footerClass)
    
    // Apply column styles
    if col.Width != "" {
        footerElem = footerElem.Style("width", col.Width)
    }
    if col.MinWidth != "" {
        footerElem = footerElem.Style("min-width", col.MinWidth)
    }
    if col.MaxWidth != "" {
        footerElem = footerElem.Style("max-width", col.MaxWidth)
    }
    
    return footerElem.Body(footerContent)
}

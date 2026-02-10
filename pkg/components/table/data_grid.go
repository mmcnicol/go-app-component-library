// pkg/components/table/data_grid.go
package table

import (
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// DataGridProps defines properties for the advanced data grid
type DataGridProps struct {
    TableProps
    // Selection
    Selectable      bool
    MultiSelect     bool
    SelectedRows    map[string]bool
    OnSelectionChange func(selectedRows map[string]interface{})
    
    // Pagination
    Pagination      bool
    PageSize        int
    CurrentPage     int
    TotalItems      int
    OnPageChange    func(page int, pageSize int)
    
    // Filtering
    Filters         map[string]string
    OnFilterChange  func(filters map[string]string)
    
    // Virtualization
    VirtualScroll   bool
    RowHeight       int
    VisibleRows     int
    
    // Additional features
    ResizableColumns bool
    ReorderableColumns bool
    ColumnVisibility map[string]bool
    Actions         []GridAction
}

// GridAction defines actions that can be performed on the grid
type GridAction struct {
    ID       string
    Label    string
    Icon     string
    Handler  func(ctx app.Context, selectedRows []map[string]interface{})
    Disabled bool
}

// DataGrid represents an advanced data grid component
type DataGrid struct {
    app.Compo
    props           DataGridProps
    internalState   GridState
}

// GridState holds the internal state of the data grid
type GridState struct {
    SelectedRows    map[string]bool
    CurrentPage     int
    Filters         map[string]string
    SortBy          string
    SortOrder       string
    ColumnWidths    map[string]int
    ColumnOrder     []string
    VisibleColumns  map[string]bool
}

// Render returns the UI representation of the data grid
func (d *DataGrid) Render() app.UI {
    return app.Div().
        Class("data-grid").
        Body(
            d.renderToolbar(),
            d.renderGrid(),
            d.renderFooter(),
        )
}

func (d *DataGrid) renderToolbar() app.UI {
    if !d.props.Selectable && len(d.props.Actions) == 0 && !d.props.Pagination {
        return app.Div()
    }
    
    return app.Div().
        Class("data-grid__toolbar").
        Body(
            d.renderSelectionInfo(),
            d.renderActions(),
            d.renderPaginationControls(),
        )
}

func (d *DataGrid) renderSelectionInfo() app.UI {
    if !d.props.Selectable || len(d.internalState.SelectedRows) == 0 {
        return app.Div()
    }
    
    selectedCount := len(d.internalState.SelectedRows)
    return app.Div().
        Class("data-grid__selection-info").
        Body(
            app.Textf("%d item%s selected", selectedCount, func() string {
                if selectedCount == 1 {
                    return ""
                }
                return "s"
            }()),
            app.Button().
                Type("button").
                Class("btn btn--text btn--small").
                OnClick(d.clearSelection).
                Text("Clear selection"),
        )
}

func (d *DataGrid) renderActions() app.UI {
    if len(d.props.Actions) == 0 {
        return app.Div()
    }
    
    i := &icon.Icon{} // Use your library icon component
    var actions []app.UI
    for _, action := range d.props.Actions {
        // ... logic for disabled status ...
        
        actions = append(actions, app.Button().
            Type("button").
            Class("btn btn--secondary btn--small").
            Disabled(disabled).
            OnClick(d.handleAction(action)).
            Body(
                i.GetIcon(action.Icon, 16), // Fix: convert string to UI
                app.Text(action.Label),
            ))
    }
    return app.Div().Class("data-grid__actions").Body(actions...)
}

func (d *DataGrid) renderGrid() app.UI {
    // Calculate visible columns
    var visibleColumns []Column
    for _, col := range d.props.Columns {
        if visible, ok := d.internalState.VisibleColumns[col.ID]; !ok || visible {
            visibleColumns = append(visibleColumns, col)
        }
    }
    
    // Create table props with selection column if needed
    tableProps := d.props.TableProps
    tableProps.Columns = visibleColumns
    
    if d.props.Selectable {
        // Add selection column as first column
        selectionColumn := Column{
            ID:       "selection",
            Header:   d.renderSelectionHeader(),
            Width:    "40px",
            Align:    "center",
            CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
                rowData := d.props.Data[rowIndex]
                rowKey := d.getRowKey(rowData, rowIndex)
                isSelected := d.internalState.SelectedRows[rowKey]
                
                return app.Input().
                    Type("checkbox").
                    Checked(isSelected).
                    OnChange(d.handleRowSelection(rowKey, rowData)).
                    Class("data-grid__checkbox")
            },
        }
        tableProps.Columns = append([]Column{selectionColumn}, tableProps.Columns...)
    }
    
    // Handle row click for selection if single select
    if d.props.Selectable && !d.props.MultiSelect {
        originalOnRowClick := tableProps.OnRowClick
        tableProps.OnRowClick = func(ctx app.Context, rowData map[string]interface{}, index int) {
            rowKey := d.getRowKey(rowData, index)
            d.toggleRowSelection(rowKey, rowData)
            if originalOnRowClick != nil {
                originalOnRowClick(ctx, rowData, index)
            }
        }
    }
    
    // Use sortable table if needed
    var table app.UI
    if len(d.props.Columns) > 0 {
        for _, col := range d.props.Columns {
            if col.Sortable {
                table = &SortableTable{
                    props: SortableTableProps{
                        TableProps: tableProps,
                        InitialSortBy: d.internalState.SortBy,
                        InitialSortOrder: d.internalState.SortOrder,
                        OnSortChange: d.handleSortChange,
                    },
                }
                break
            }
        }
    }
    
    if table == nil {
        table = &Table{props: tableProps}
    }
    
    return table
}

func (d *DataGrid) renderSelectionHeader() app.UI {
    if !d.props.MultiSelect {
        return app.Text("Select")
    }
    
    allSelected := len(d.internalState.SelectedRows) == len(d.props.Data)
    indeterminate := len(d.internalState.SelectedRows) > 0 && !allSelected
    
    return app.Input().
        Type("checkbox").
        Checked(allSelected).
        Indeterminate(indeterminate).
        OnChange(d.handleSelectAll).
        Class("data-grid__checkbox data-grid__checkbox--header")
}

func (d *DataGrid) renderFooter() app.UI {
    if !d.props.Pagination {
        return app.Div()
    }
    
    totalPages := (d.props.TotalItems + d.props.PageSize - 1) / d.props.PageSize
    
    return app.Div().
        Class("data-grid__footer").
        Body(
            d.renderPaginationInfo(totalPages),
            d.renderPageNavigation(totalPages),
            d.renderPageSizeSelector(),
        )
}

func (d *DataGrid) renderPaginationInfo(totalPages int) app.UI {
    start := (d.internalState.CurrentPage-1)*d.props.PageSize + 1
    end := min(d.internalState.CurrentPage*d.props.PageSize, d.props.TotalItems)
    
    return app.Div().
        Class("data-grid__pagination-info").
        Textf("Showing %d-%d of %d items", start, end, d.props.TotalItems)
}

func (d *DataGrid) renderPageNavigation(totalPages int) app.UI {
    return app.Div().
        Class("data-grid__page-navigation").
        Body(
            app.Button().
                Type("button").
                Class("btn btn--icon btn--small").
                Disabled(d.internalState.CurrentPage == 1).
                OnClick(d.goToFirstPage).
                Body(app.Text("⏮")),
            app.Button().
                Type("button").
                Class("btn btn--icon btn--small").
                Disabled(d.internalState.CurrentPage == 1).
                OnClick(d.goToPreviousPage).
                Body(app.Text("◀")),
            app.Input().
                Type("number").
                Value(d.internalState.CurrentPage).
                Min(1).
                Max(totalPages).
                OnChange(d.handlePageInput).
                Class("data-grid__page-input"),
            app.Text("/"),
            app.Text(totalPages),
            app.Button().
                Type("button").
                Class("btn btn--icon btn--small").
                Disabled(d.internalState.CurrentPage == totalPages).
                OnClick(d.goToNextPage).
                Body(app.Text("▶")),
            app.Button().
                Type("button").
                Class("btn btn--icon btn--small").
                Disabled(d.internalState.CurrentPage == totalPages).
                OnClick(d.goToLastPage).
                Body(app.Text("⏭")),
        )
}

func (d *DataGrid) renderPageSizeSelector() app.UI {
    pageSizes := []int{10, 25, 50, 100}
    
    var options []app.UI
    for _, size := range pageSizes {
        options = append(options, app.Option().
            Value(size).
            Text(string(size)).
            Selected(size == d.props.PageSize))
    }
    
    return app.Select().
        Value(d.props.PageSize).
        OnChange(d.handlePageSizeChange).
        Class("data-grid__page-size").
        Body(options...)
}

// Helper functions for pagination
func (d *DataGrid) goToFirstPage(ctx app.Context, e app.Event) {
    d.internalState.CurrentPage = 1
    ctx.Update()
    d.notifyPageChange()
}

func (d *DataGrid) goToPreviousPage(ctx app.Context, e app.Event) {
    if d.internalState.CurrentPage > 1 {
        d.internalState.CurrentPage--
        ctx.Update()
        d.notifyPageChange()
    }
}

func (d *DataGrid) goToNextPage(ctx app.Context, e app.Event) {
    totalPages := (d.props.TotalItems + d.props.PageSize - 1) / d.props.PageSize
    if d.internalState.CurrentPage < totalPages {
        d.internalState.CurrentPage++
        ctx.Update()
        d.notifyPageChange()
    }
}

func (d *DataGrid) goToLastPage(ctx app.Context, e app.Event) {
    totalPages := (d.props.TotalItems + d.props.PageSize - 1) / d.props.PageSize
    d.internalState.CurrentPage = totalPages
    ctx.Update()
    d.notifyPageChange()
}

func (d *DataGrid) handlePageInput(ctx app.Context, e app.Event) {
    // Implementation for page input handling
}

func (d *DataGrid) handlePageSizeChange(ctx app.Context, e app.Event) {
    // Implementation for page size change
}

func (d *DataGrid) notifyPageChange() {
    if d.props.OnPageChange != nil {
        d.props.OnPageChange(d.internalState.CurrentPage, d.props.PageSize)
    }
}

func (d *DataGrid) getRowKey(rowData map[string]interface{}, index int) string {
	// If the data has an "ID", use it; otherwise, fallback to the row index
	if id, ok := rowData["ID"].(string); ok {
		return id
	}
	return fmt.Sprintf("row-%d", index)
}

func (d *DataGrid) toggleRowSelection(key string, data map[string]interface{}) {
	if d.internalState.SelectedRows == nil {
		d.internalState.SelectedRows = make(map[string]bool)
	}
	d.internalState.SelectedRows[key] = !d.internalState.SelectedRows[key]
}

func (d *DataGrid) clearSelection(ctx app.Context, e app.Event) {
	d.internalState.SelectedRows = make(map[string]bool)
	ctx.Update()
}

func (d *DataGrid) handleRowSelection(key string, data map[string]interface{}) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		d.toggleRowSelection(key, data)
		ctx.Update()
	}
}

func (d *DataGrid) handleAction(action GridAction) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		if action.Handler != nil {
			// Convert internal map to slice for the handler
			var selected []map[string]interface{}
			// Note: This requires tracking the actual data objects in internalState 
			// if you need to pass full objects back.
			action.Handler(ctx, selected)
		}
	}
}

func (d *DataGrid) handleSortChange(sortBy string, sortOrder string) {
	d.internalState.SortBy = sortBy
	d.internalState.SortOrder = sortOrder
}

func (d *DataGrid) renderPaginationControls() app.UI {
    if !d.props.Pagination {
        return app.Nil
    }
    // You can return the same nav used in the footer
    totalPages := (d.props.TotalItems + d.props.PageSize - 1) / d.props.PageSize
    return d.renderPageNavigation(totalPages)
}

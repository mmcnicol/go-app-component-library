// pkg/components/table/sortable_table.go
package table

import (
    "sort"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// SortableTableProps extends TableProps with sorting capabilities
type SortableTableProps struct {
    TableProps
    InitialSortBy    string
    InitialSortOrder string // "asc", "desc"
    OnSortChange     func(sortBy string, sortOrder string)
}

// SortableTable adds sorting functionality to the base table
type SortableTable struct {
    app.Compo
    props      SortableTableProps
    sortBy     string
    sortOrder  string
    sortedData []map[string]interface{}
}

// OnMount initializes the component
func (s *SortableTable) OnMount(ctx app.Context) {
    s.sortBy = s.props.InitialSortBy
    s.sortOrder = s.props.InitialSortOrder
    s.sortedData = make([]map[string]interface{}, len(s.props.Data))
    copy(s.sortedData, s.props.Data)
    s.sortData()
}

// Render returns the UI representation of the sortable table
func (s *SortableTable) Render() app.UI {
    // Update table props with sorted data and custom header renderer
    tableProps := s.props.TableProps
    tableProps.Data = s.sortedData
    
    // Replace header renderer to add sorting indicators
    for i := range tableProps.Columns {
        col := &tableProps.Columns[i]
        if col.Sortable {
            originalRenderer := col.HeaderRenderer
            col.HeaderRenderer = func(col Column, colIndex int) app.UI {
                var content []app.UI
                
                // Original content
                if originalRenderer != nil {
                    content = append(content, originalRenderer(col, colIndex))
                } else {
                    content = append(content, app.Text(col.Header))
                }
                
                // Add sorting indicator
                if s.sortBy == col.Accessor {
                    icon := "↑"
                    if s.sortOrder == "desc" {
                        icon = "↓"
                    }
                    content = append(content, app.Span().
                        Class("table__sort-indicator").
                        Text(icon))
                } else {
                    content = append(content, app.Span().
                        Class("table__sort-indicator table__sort-indicator--inactive").
                        Text("↕"))
                }
                
                return app.Button().
                    Type("button").
                    Class("table__sort-button").
                    OnClick(s.handleSortClick(col.Accessor)).
                    Body(content...)
            }
        }
    }
    
    return &Table{props: tableProps}
}

func (s *SortableTable) handleSortClick(accessor string) func(ctx app.Context, e app.Event) {
    return func(ctx app.Context, e app.Event) {
        e.PreventDefault()
        
        if s.sortBy == accessor {
            // Toggle order if clicking the same column
            if s.sortOrder == "asc" {
                s.sortOrder = "desc"
            } else {
                s.sortOrder = "asc"
            }
        } else {
            // Set new column and default to ascending
            s.sortBy = accessor
            s.sortOrder = "asc"
        }
        
        s.sortData()
        ctx.Update()
        
        if s.props.OnSortChange != nil {
            s.props.OnSortChange(s.sortBy, s.sortOrder)
        }
    }
}

func (s *SortableTable) sortData() {
    if s.sortBy == "" || len(s.sortedData) == 0 {
        return
    }
    
    sort.Slice(s.sortedData, func(i, j int) bool {
        valI := s.sortedData[i][s.sortBy]
        valJ := s.sortedData[j][s.sortBy]
        
        // Handle different data types
        switch vI := valI.(type) {
        case string:
            vJ := valJ.(string)
            if s.sortOrder == "asc" {
                return vI < vJ
            }
            return vI > vJ
            
        case int:
            vJ := valJ.(int)
            if s.sortOrder == "asc" {
                return vI < vJ
            }
            return vI > vJ
            
        case float64:
            vJ := valJ.(float64)
            if s.sortOrder == "asc" {
                return vI < vJ
            }
            return vI > vJ
            
        default:
            // For other types, convert to string
            strI := fmt.Sprintf("%v", vI)
            strJ := fmt.Sprintf("%v", valJ)
            if s.sortOrder == "asc" {
                return strI < strJ
            }
            return strI > strJ
        }
    })
}

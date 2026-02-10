//go:build dev
// pkg/components/table/table_stories.go
package table

import (
	"fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {
    storybook.Register("Data", "Basic Table",
        map[string]*storybook.Control{
            "Striped":      {Label: "Striped", Type: storybook.ControlBool, Value: true},
            "Bordered":     {Label: "Bordered", Type: storybook.ControlBool, Value: false},
            "Hoverable":    {Label: "Hoverable", Type: storybook.ControlBool, Value: true},
            "Compact":      {Label: "Compact", Type: storybook.ControlBool, Value: false},
            "Responsive":   {Label: "Responsive", Type: storybook.ControlBool, Value: true},
            "Loading":      {Label: "Loading", Type: storybook.ControlBool, Value: false},
            "Empty":        {Label: "Show Empty State", Type: storybook.ControlBool, Value: false},
            "DataSize":     {Label: "Data Size", Type: storybook.ControlSelect, Value: "5", Options: []string{"0", "3", "5", "10", "20"}},
        },
        func(controls map[string]*storybook.Control) app.UI {
            striped := controls["Striped"].Value.(bool)
            bordered := controls["Bordered"].Value.(bool)
            hoverable := controls["Hoverable"].Value.(bool)
            compact := controls["Compact"].Value.(bool)
            responsive := controls["Responsive"].Value.(bool)
            loading := controls["Loading"].Value.(bool)
            empty := controls["Empty"].Value.(bool)
            dataSize := controls["DataSize"].Value.(string)

            // Sample data
            columns := []Column{
                {
                    ID:       "id",
                    Header:   "ID",
                    Accessor: "id",
                    Width:    "80px",
                    Sortable: true,
                },
                {
                    ID:       "name",
                    Header:   "Name",
                    Accessor: "name",
                    Width:    "200px",
                    Sortable: true,
                },
                {
                    ID:       "email",
                    Header:   "Email",
                    Accessor: "email",
                    Width:    "250px",
                },
                {
                    ID:       "role",
                    Header:   "Role",
                    Accessor: "role",
                    Width:    "120px",
                    Align:    "center",
                },
                {
                    ID:       "status",
                    Header:   "Status",
                    Accessor: "status",
                    Width:    "100px",
                    Align:    "center",
                    CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
                        status := data.(string)
                        badgeClass := "badge "
                        if status == "Active" {
                            badgeClass += "badge--success"
                        } else if status == "Pending" {
                            badgeClass += "badge--warning"
                        } else {
                            badgeClass += "badge--error"
                        }
                        return app.Span().Class(badgeClass).Text(status)
                    },
                },
                {
                    ID:       "actions",
                    Header:   "Actions",
                    Accessor: "actions",
                    Width:    "150px",
                    Align:    "center",
                    CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
                        return app.Div().Class("btn-group").Body(
                            app.Button().Class("btn btn--small btn--text").Text("Edit"),
                            app.Button().Class("btn btn--small btn--text btn--danger").Text("Delete"),
                        )
                    },
                },
            }

            // Generate sample data based on dataSize
            var data []map[string]interface{}
            if empty {
                dataSize = "0"
            }
            
            size := 0
            switch dataSize {
            case "0":
                size = 0
            case "3":
                size = 3
            case "5":
                size = 5
            case "10":
                size = 10
            case "20":
                size = 20
            }

            for i := 1; i <= size; i++ {
                status := "Active"
                if i%3 == 0 {
                    status = "Pending"
                } else if i%5 == 0 {
                    status = "Inactive"
                }
                
                role := "User"
                if i%4 == 0 {
                    role = "Admin"
                } else if i%7 == 0 {
                    role = "Moderator"
                }

                data = append(data, map[string]interface{}{
                    "id":     fmt.Sprintf("USR%04d", i),
                    "name":   fmt.Sprintf("User %d", i),
                    "email":  fmt.Sprintf("user%d@example.com", i),
                    "role":   role,
                    "status": status,
                    "actions": nil,
                })
            }

            return &Table{
                props: TableProps{
                    Columns:    columns,
                    Data:       data,
                    Striped:    striped,
                    Bordered:   bordered,
                    Hoverable:  hoverable,
                    Compact:    compact,
                    Responsive: responsive,
                    Loading:    loading,
                    RowKey:     "id",
                    OnRowClick: func(ctx app.Context, rowData map[string]interface{}, index int) {
                        app.Logf("Row clicked: %v", rowData["id"])
                    },
                    EmptyState: app.Div().Class("text-center p-8").Body(
                        app.Div().Class("text-gray-400 mb-4").Text("📊"),
                        app.H3().Class("text-lg font-semibold mb-2").Text("No Data Available"),
                        app.P().Class("text-gray-600").Text("There are no records to display."),
                    ),
                    LoadingState: app.Div().Class("text-center p-8").Body(
                        app.Div().Class("spinner mx-auto mb-4"),
                        app.P().Class("text-gray-600").Text("Loading table data..."),
                    ),
                },
            }
        },
    )

    storybook.Register("Data", "Sortable Table",
        map[string]*storybook.Control{
            "SortBy":       {Label: "Sort By", Type: storybook.ControlSelect, Value: "name", Options: []string{"id", "name", "department", "salary", "hireDate"}},
            "SortOrder":    {Label: "Sort Order", Type: storybook.ControlSelect, Value: "asc", Options: []string{"asc", "desc"}},
            "DataSize":     {Label: "Data Size", Type: storybook.ControlSelect, Value: "10", Options: []string{"5", "10", "20", "50"}},
        },
        func(controls map[string]*storybook.Control) app.UI {
            sortBy := controls["SortBy"].Value.(string)
            sortOrder := controls["SortOrder"].Value.(string)
            dataSize := controls["DataSize"].Value.(string)

			columns := []Column{
				{
					ID:       "id",
					Header:   "ID",
					Accessor: "id",
					Width:    "100px",
					Sortable: true,
				},
				{
					ID:       "name",
					Header:   "Name",
					Accessor: "name",
					Width:    "200px",
					Sortable: true,
				},
				{
					ID:       "department",
					Header:   "Department",
					Accessor: "department",
					Width:    "150px",
					Sortable: true,
				},
				{
					ID:       "salary",
					Header:   "Salary",
					Accessor: "salary",
					Width:    "120px",
					Align:    "right",
					Sortable: true,
					CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
						salary := data.(float64)
						return app.Text(fmt.Sprintf("$%.2f", salary))
					},
				},
				{
					ID:       "hireDate",
					Header:   "Hire Date",
					Accessor: "hireDate",
					Width:    "120px",
					Sortable: true,
				},
			}

			// Generate sample data
			size := 10
			switch dataSize {
			case "5":
				size = 5
			case "10":
				size = 10
			case "20":
				size = 20
			case "50":
				size = 50
			}

			departments := []string{"Engineering", "Marketing", "Sales", "HR", "Finance"}
			data := make([]map[string]interface{}, size)
			
			for i := 0; i < size; i++ {
				deptIndex := i % len(departments)
				data[i] = map[string]interface{}{
					"id":         fmt.Sprintf("EMP%03d", i+1),
					"name":       fmt.Sprintf("Employee %d", i+1),
					"department": departments[deptIndex],
					"salary":     50000.0 + float64(i)*1000.0,
					"hireDate":   fmt.Sprintf("2023-%02d-%02d", (i%12)+1, (i%28)+1),
				}
			}

			// Create the SortableTable component
			sortableTable := &SortableTable{
				props: SortableTableProps{
					TableProps: TableProps{
						Columns:   columns,
						Data:      data, // Pass the data here
						Striped:   true,
						Hoverable: true,
						RowKey:    "id",
					},
					InitialSortBy:    initialSort,
					InitialSortOrder: sortOrder,
					OnSortChange: func(sortBy string, sortOrder string) {
						app.Logf("Sort changed: %s %s", sortBy, sortOrder)
					},
				},
			}

			// Initialize the component
			sortableTable.sortBy = initialSort
			sortableTable.sortOrder = sortOrder
			sortableTable.sortedData = make([]map[string]interface{}, len(data))
			copy(sortableTable.sortedData, data)
			if initialSort != "" {
				// Sort the data initially
				sortableTable.sortData()
			}

			return &SortableTable{
                props: SortableTableProps{
                    TableProps: TableProps{
                        Columns:   columns,
                        Data:      data,
                        Striped:   true,
                        Hoverable: true,
                        RowKey:    "id",
                    },
                    InitialSortBy:    sortBy,
                    InitialSortOrder: sortOrder,
                    OnSortChange: func(newSortBy string, newSortOrder string) {
                        // Update the controls
                        controls["SortBy"].Value = newSortBy
                        controls["SortOrder"].Value = newSortOrder
                        
                        // The Shell will detect control changes and re-render
                    },
                },
            }
        },
    )

    storybook.Register("Data", "Data Grid",
        map[string]*storybook.Control{
            "Selectable":   {Label: "Selectable", Type: storybook.ControlBool, Value: true},
            "MultiSelect":  {Label: "Multi Select", Type: storybook.ControlBool, Value: true},
            "Pagination":   {Label: "Pagination", Type: storybook.ControlBool, Value: true},
            "PageSize":     {Label: "Page Size", Type: storybook.ControlSelect, Value: "10", Options: []string{"5", "10", "25", "50"}},
            "DataSize":     {Label: "Total Items", Type: storybook.ControlSelect, Value: "45", Options: []string{"15", "30", "45", "100"}},
        },
        func(controls map[string]*storybook.Control) app.UI {
            selectable := controls["Selectable"].Value.(bool)
            multiSelect := controls["MultiSelect"].Value.(bool)
            pagination := controls["Pagination"].Value.(bool)
            pageSize := controls["PageSize"].Value.(string)
            dataSize := controls["DataSize"].Value.(string)

            columns := []Column{
                {
                    ID:       "product",
                    Header:   "Product",
                    Accessor: "product",
                    Width:    "200px",
                },
                {
                    ID:       "category",
                    Header:   "Category",
                    Accessor: "category",
                    Width:    "150px",
                },
                {
                    ID:       "price",
                    Header:   "Price",
                    Accessor: "price",
                    Width:    "100px",
                    Align:    "right",
                    CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
                        price := data.(float64)
                        return app.Text(fmt.Sprintf("$%.2f", price))
                    },
                },
                {
                    ID:       "stock",
                    Header:   "Stock",
                    Accessor: "stock",
                    Width:    "100px",
                    Align:    "center",
                    CellRenderer: func(data interface{}, rowIndex int, colIndex int) app.UI {
                        stock := data.(int)
                        if stock > 20 {
                            return app.Span().Class("text-green-600").Text(fmt.Sprintf("%d", stock))
                        } else if stock > 0 {
                            return app.Span().Class("text-yellow-600").Text(fmt.Sprintf("%d", stock))
                        }
                        return app.Span().Class("text-red-600").Text(fmt.Sprintf("%d", stock))
                    },
                },
                {
                    ID:       "lastUpdated",
                    Header:   "Last Updated",
                    Accessor: "lastUpdated",
                    Width:    "150px",
                },
            }

            // Parse values
            pageSizeInt := 10
            switch pageSize {
            case "5":
                pageSizeInt = 5
            case "10":
                pageSizeInt = 10
            case "25":
                pageSizeInt = 25
            case "50":
                pageSizeInt = 50
            }

            totalItems := 45
            switch dataSize {
            case "15":
                totalItems = 15
            case "30":
                totalItems = 30
            case "45":
                totalItems = 45
            case "100":
                totalItems = 100
            }

            // Generate sample data (just for current page)
            categories := []string{"Electronics", "Clothing", "Books", "Home & Garden", "Toys"}
            products := []string{"Laptop", "Smartphone", "Tablet", "Headphones", "Monitor", "Keyboard", "Mouse"}
            
            var data []map[string]interface{}
            for i := 0; i < pageSizeInt && i < totalItems; i++ {
                productIndex := i % len(products)
                categoryIndex := i % len(categories)
                
                data = append(data, map[string]interface{}{
                    "product":     fmt.Sprintf("%s Pro Max", products[productIndex]),
                    "category":    categories[categoryIndex],
                    "price":       199.99 + float64(i)*50.0,
                    "stock":       50 - i%30,
                    "lastUpdated": fmt.Sprintf("2024-01-%02d", (i%28)+1),
                })
            }

            // Define actions for the grid
            actions := []GridAction{
                {
                    ID:    "edit",
                    Label: "Edit",
                    Icon:  "edit",
                    Handler: func(ctx app.Context, selectedRows []map[string]interface{}) {
                        app.Logf("Edit action triggered on %d rows", len(selectedRows))
                    },
                },
                {
                    ID:    "delete",
                    Label: "Delete",
                    Icon:  "trash",
                    Handler: func(ctx app.Context, selectedRows []map[string]interface{}) {
                        app.Logf("Delete action triggered on %d rows", len(selectedRows))
                    },
                    Disabled: false,
                },
                {
                    ID:    "export",
                    Label: "Export",
                    Icon:  "download",
                    Handler: func(ctx app.Context, selectedRows []map[string]interface{}) {
                        app.Logf("Export action triggered on %d rows", len(selectedRows))
                    },
                },
            }

            return &DataGrid{
                props: DataGridProps{
                    TableProps: TableProps{
                        Columns:   columns,
                        Data:      data,
                        Striped:   true,
                        Hoverable: true,
                    },
                    Selectable:   selectable,
                    MultiSelect:  multiSelect,
                    Pagination:   pagination,
                    PageSize:     pageSizeInt,
                    CurrentPage:  1,
                    TotalItems:   totalItems,
                    Actions:      actions,
                    OnPageChange: func(page int, pageSize int) {
                        app.Logf("Page changed to %d with size %d", page, pageSize)
                    },
                    OnSelectionChange: func(selectedRows map[string]interface{}) {
                        app.Logf("Selection changed: %d rows selected", len(selectedRows))
                    },
                },
            }
        },
    )
}

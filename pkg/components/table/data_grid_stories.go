//go:build dev
// pkg/components/table/table_stories.go
package table

import (
	"fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {

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

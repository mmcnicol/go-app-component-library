//go:build dev
// pkg/components/table/table_stories.go
package table

import (
	"fmt"
    "github.com/maxence-charriere/go-app/v10/pkg/app"
    "github.com/mmcnicol/go-app-component-library/pkg/storybook"
)

func init() {

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
    
}

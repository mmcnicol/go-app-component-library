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
    
}

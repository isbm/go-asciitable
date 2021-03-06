package asciitable

import (
	"fmt"
	"reflect"
	"strings"
)

type TableData struct {
	header []string
	data   [][]string
}

/*
NewTableData object constructor
*/
func NewTableData() *TableData {
	tableData := new(TableData)
	tableData.data = make([][]string, 0)
	return tableData
}

func (tableData *TableData) SetData(data [][]interface{}) *TableData {
	if len(data) < 1 && len(data[0]) < 1 {
		panic("ERROR: Data should to be two-dimentional array with at least one row, containing at least one element.")
	}

	for _, row := range data {
		tableData.AddRow(row...)
	}

	return tableData
}

func (tableData *TableData) AddRow(row ...interface{}) *TableData {
	data := make([]string, len(row))
	for idx, rowData := range row {
		var cellData string
		switch reflect.TypeOf(rowData).Kind() {
		case reflect.Int:
			cellData = fmt.Sprintf("%d", rowData)
		case reflect.Slice, reflect.Array, reflect.Map:
			cellData = fmt.Sprintf("*### Error: %s ###*", reflect.TypeOf(rowData))
		case reflect.String:
			cellData = strings.TrimSpace(rowData.(string))
		default:
			cellData = fmt.Sprintf("%v", rowData)
		}
		data[idx] = cellData
	}

	if len(row) > 0 {
		tableData.data = append(tableData.data, data)
	}

	return tableData
}

// Get raw table data
func (tableData *TableData) GetData() *[][]string {
	return &tableData.data
}

// Get table header data
func (tableData *TableData) GetHeader() *[]string {
	return &tableData.header
}

// Get number of columns in the table
func (tableData *TableData) GetColsNum() int {
	cols := 0
	if tableData.GetRowsNum() > 0 {
		cols = len(tableData.data[0])
	}
	return cols
}

// Get number of rows in the table
func (tableData *TableData) GetRowsNum() int {
	return len(tableData.data)
}

/*
Set Header of the table. Each string represents a column name.
Previous data is wept away.
*/
func (tableData *TableData) SetHeader(titles ...string) *TableData {
	tableData.header = make([]string, len(titles))
	copy(tableData.header, titles)

	return tableData
}

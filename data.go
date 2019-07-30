package asciitable

import (
	"fmt"
	"reflect"
)

type TableData struct {
	data [][]string
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

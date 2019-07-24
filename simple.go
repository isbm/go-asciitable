package asciitable

import (
	"fmt"
	"strconv"
	"strings"
)

// Border mode. Internal use.
const (
	_borderTop = iota
	_borderInner
	_borderBottom
)

// Table object
type simpleTable struct {
	header      []string
	rowsData    [][]string
	rowsCount   uint64
	headerAlign int
	cellAlign   int
	border      byte
	style       borderStyle
}

/*
NewSimpleTable object constructor
*/
func NewSimpleTable(style *borderStyle) *simpleTable {
	table := new(simpleTable)
	table.rowsData = make([][]string, 0)
	table.rowsCount = 0
	table.headerAlign = ALIGN_LEFT
	table.cellAlign = ALIGN_LEFT
	if style == nil {
		style = NewBorderStyle(-1, -1)
	}
	table.style = *style

	return table
}

/*
Set Header of the table. Each string represents a column name.
Previous data is wept away.
*/
func (table *simpleTable) Header(titles ...string) {
	table.header = make([]string, len(titles))
	copy(table.header, titles)
}

/*
Add row to the data. Mostly strings... :-)
*/
func (table *simpleTable) AddRow(data ...interface{}) {
	row := make([]string, len(data))
	for idx, rowData := range data {
		var cellData string
		switch rowData.(type) {
		case int:
			cellData = fmt.Sprintf("%d", rowData) // Formatting later
		default:
			cellData = fmt.Sprintf("%v", rowData)
		}
		row[idx] = cellData
	}

	if len(row) > 0 {
		table.rowsData = append(table.rowsData, row)
		table.rowsCount++
	}
}

// Calculate row widths for maximum widest data
func (table *simpleTable) getRowWidths() []int {
	widths := make([]int, len(table.header))
	for idx := range widths {
		widths[idx] = 0
	}

	for idx, title := range table.header {
		if len(title) > widths[idx] {
			widths[idx] = len(title)
		}
	}

	for _, rData := range table.rowsData {
		for idx, data := range rData {
			if len(data) > widths[idx] {
				widths[idx] = len(data)
			}
		}
	}

	return widths
}

func (table *simpleTable) renderCell(data string, width int, first bool) string {
	w := strconv.Itoa(width)
	var cell string
	switch table.headerAlign {
	case ALIGN_RIGHT:
		cell = fmt.Sprintf("%"+w+"s", data)
	case ALIGN_CENTER:
		cell = fmt.Sprintf("%-"+w+"s", fmt.Sprintf("%"+w+"s", data))
	default:
		cell = fmt.Sprintf("%-"+w+"s", data)
	}
	return cell
}

/*
Renders border, but not exactly a *border* but row between the table data,
which is typically either top of the header (outer border of the table),
row under the header, row between the regular cells or bottom row (outer
border of the table, again).
*/
func (table *simpleTable) renderBorder(borderType int) string {
	rowWidths := table.getRowWidths()
	var border string
	switch borderType {
	case _borderTop:
		border = table.style.outer.LeftTop()
		for idx, width := range rowWidths {
			border += strings.Repeat(table.style.outer.HorisontalLine(), width)
			if idx < len(rowWidths)-1 {
				border += table.style.inner.CenterTop()
			} else {
				border += table.style.outer.RightTop()
			}
		}
	case _borderBottom:
		border = table.style.outer.LeftBottom()
		for idx, width := range rowWidths {
			border += strings.Repeat(table.style.outer.HorisontalLine(), width)
			if idx < len(rowWidths)-1 {
				border += table.style.inner.CenterBottom()
			} else {
				border += table.style.outer.RightBottom()
			}
		}
	case _borderInner:
		border = table.style.inner.LeftMiddle()
		for idx, width := range rowWidths {
			border += strings.Repeat(table.style.inner.HorisontalLine(), width)
			if idx < len(rowWidths)-1 {
				border += table.style.inner.CenterMiddle()
			} else {
				border += table.style.inner.RightMiddle()
			}
		}
	}
	return border
}

// Takes padded cells data and renders to the row
func (table *simpleTable) renderRow(cells []string) string {
	rowWidths := table.getRowWidths()
	var row string
	for idx, cell := range cells {
		if idx < 1 {
			row += table.style.inner.VerticalLine()
		}
		row += table.renderCell(cell, rowWidths[idx], idx == 0)
		row += table.style.inner.VerticalLine()
	}
	return row
}

// Renders table as a string
func (table *simpleTable) Render() string {
	render := make([]string, 0)

	if len(table.header) > 0 {
		render = append(render, []string{
			table.renderBorder(_borderTop),
			table.renderRow(table.header),
			table.renderBorder(_borderInner),
		}...)
	}

	for _, row := range table.rowsData[:len(table.rowsData)-1] {
		render = append(render, []string{
			table.renderRow(row),
			table.renderBorder(_borderInner),
		}...)
	}
	render = append(render, []string{
		table.renderRow(table.rowsData[len(table.rowsData)-1]),
		table.renderBorder(_borderBottom),
	}...)

	return strings.Join(render, "\n")
}

package asciitable

import (
	"github.com/isbm/textwrap"
	"regexp"
	"strings"
)

// Border mode. Internal use.
const (
	_borderTop = iota
	_borderInner
	_borderBottom
	_borderHeader
	_ansiRegex = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
)

// Table object
type simpleTable struct {
	rowsData        *TableData
	rowsCount       uint64
	headerAlign     int
	columnsAlign    []int
	columnsTextWrap []bool
	style           *borderStyle
	widthTable      int
	widthColumns    []int
	widthData       int
	padding         int
	wrapText        bool
	stripAnsiRegex  *regexp.Regexp
}

/*
NewSimpleTable object constructor
*/
func NewSimpleTable(data *TableData, style *borderStyle) *simpleTable {
	table := new(simpleTable)
	if data == nil {
		data = NewTableData()
	}
	table.rowsData = data
	table.rowsCount = 0
	table.headerAlign = ALIGN_LEFT

	// Set default column align and nowrap
	table.columnsAlign = make([]int, data.GetColsNum())
	table.columnsTextWrap = make([]bool, len(table.columnsAlign))
	for idx := range table.columnsAlign {
		table.columnsAlign[idx] = ALIGN_LEFT
		table.columnsTextWrap[idx] = true
	}

	// Set style
	if style == nil {
		style = NewBorderStyle(-1, -1)
	}
	table.style = style

	table.widthTable, _ = GetTerminalSize()
	table.widthColumns = make([]int, 0)
	table.padding = 0
	table.wrapText = false
	table.stripAnsiRegex = regexp.MustCompile(_ansiRegex)

	return table
}

// SetWrapText wraps text in all cells instead of trimming it to the max width.
func (table *simpleTable) SetTextWrap(wrap bool) *simpleTable {
	table.wrapText = wrap
	return table
}

// Set cell padding
func (table *simpleTable) SetCellPadding(width int) *simpleTable {
	table.padding = width
	return table
}

/*
Set overall table width (chars)
*/
func (table *simpleTable) SetWidth(width int) *simpleTable {
	table.widthTable = width
	return table
}

/*
Set column width (chars). If columns contains only one value and it is -1,
then width applies to all columns at once.
*/
func (table *simpleTable) SetColWidth(width int, columns ...int) *simpleTable {
	colsNum := table.Data().GetColsNum()
	if colsNum == 0 {
		panic("An attempt to set columns while no header or data has been set")
	} else if len(table.widthColumns) == 0 {
		table.widthColumns = make([]int, colsNum)
	} else if len(columns) > colsNum {
		panic("An attempt to set more columns widths than actually in the table")
	}

	// Set width to all cells
	if len(columns) == 1 && columns[0] == -1 {
		for idx := range table.widthColumns {
			table.widthColumns[idx] = width
		}
	} else {
		// Set only specific cells
		for _, column := range columns {
			if column < colsNum {
				table.widthColumns[column] = width
			} else {
				panic("Attempt to set width of a column that does not exist")
			}
		}
	}

	return table
}

// Set column text wrap
func (table *simpleTable) SetColTextWrap(wrap bool, columns ...int) *simpleTable {
	colsNum := table.Data().GetColsNum()

	// Set wrapping to all columns
	if len(columns) == 1 && columns[0] == -1 {
		for idx := range table.columnsTextWrap {
			table.columnsTextWrap[idx] = wrap
		}
	} else {
		// Set only specific columns
		for _, column := range columns {
			if column < colsNum {
				table.columnsTextWrap[column] = wrap
			} else {
				panic("Attempt to set text wrapping attribute to a column that does not exist")
			}
		}
	}

	return table
}

// Set column align
func (table *simpleTable) SetColAlign(align int, columns ...int) *simpleTable {
	if align != ALIGN_LEFT && align != ALIGN_RIGHT && align != ALIGN_CENTER {
		panic("An attempt to set an unknown align")
	}

	colsNum := table.Data().GetColsNum()

	if colsNum == 0 {
		panic("An attempt to set column align while no header or data has been set")
	} else if len(table.widthColumns) == 0 {
		table.widthColumns = make([]int, colsNum)
	} else if len(columns) > colsNum {
		panic("An attempt to set more column align than actually in the table")
	}

	// Set align to all cells
	if len(columns) == 1 && columns[0] == -1 {
		for idx := range table.columnsAlign {
			table.columnsAlign[idx] = align
		}
	} else {
		// Set only specific cells
		for _, column := range columns {
			if column < colsNum {
				table.columnsAlign[column] = align
			} else {
				panic("Attempt to set width of a column that does not exist")
			}
		}
	}

	return table
}

// Returns table data
func (table *simpleTable) Data() *TableData {
	return table.rowsData
}

// Allow support ANSI-colored data. If the data is not stripped out,
// all the widths will be wrongly calculated
func (table *simpleTable) stripAnsi(data string) string {
	return table.stripAnsiRegex.ReplaceAllString(data, "")
}

// Sets maximum data width. Used to decide either table is narrower
// then the terminal or not. Normally should be called after
// data bulk update, since it is quite expensive.
func (table *simpleTable) setDataMaxWidth() int {
	width := 0
	for _, row := range *table.Data().GetData() {
		rowWidth := 0
		for _, cell := range row {
			rowWidth += len(table.stripAnsi(cell))
		}
		if rowWidth > width {
			width = rowWidth
		}
	}
	table.widthData = width

	return table.widthData
}

// Calculate row widths for maximum widest data
// XXX: Should be able to cache the results and do not re-calculate it from scratch for each row and element!!
func (table *simpleTable) getRowWidths() []int {
	widths := make([]int, len(*table.Data().GetHeader()))

	for idx, title := range *table.Data().GetHeader() {
		titleLength := len(table.stripAnsi(title)) + table.padding*2
		if titleLength >= widths[idx] {
			widths[idx] = titleLength
		}
	}

	for _, rowData := range *table.Data().GetData() {
		for idx, data := range rowData {
			dataLength := len(table.stripAnsi(data)) + table.padding*2
			if dataLength > widths[idx] {
				widths[idx] = dataLength
			}
		}
	}

	// Override custom widths
	if len(table.widthColumns) == len(widths) {
		// If col width is != 0, then it is specified
		for idx, colWidth := range table.widthColumns {
			if colWidth > 0 {
				widths[idx] = colWidth
			}
		}
	}

	// Set expand table if it is set explicitly or data is bigger than the table
	if table.widthTable > table.widthData || len(table.widthColumns) > 0 {
		if table.style.widthFull {
			// Adjust the last column accordingly to the table width,
			// but only if it was not explicitly specified already
			sum := 0
			for _, width := range widths {
				sum += width
			}
			if sum < table.widthTable {
				var offset int
				if table.style.outer.IS_VISIBLE {
					offset = 4
				} else {
					offset = 2
				}
				widths[len(widths)-1] = widths[len(widths)-1] + table.widthTable - sum - offset - 2
			} else {
				lastColWidth := table.widthTable - (sum - widths[len(widths)-1]) - 3
				if lastColWidth < 4 {
					lastColWidth = 4 + table.padding*2
				}
				widths[len(widths)-1] = lastColWidth
			}
		}
	} else {
		defaultWidth := table.widthTable/table.Data().GetColsNum() - 1
		for idx := range widths {
			widths[idx] = defaultWidth
		}
		// set last column width
		lastColumnWidth := (table.widthTable - (defaultWidth * (table.Data().GetColsNum() - 1))) - table.Data().GetColsNum() + 1
		if table.style.outer.IS_VISIBLE {
			lastColumnWidth -= 2
		}
		widths[len(widths)-1] = lastColumnWidth
	}

	return widths
}

// Support ANSI escape
func (table *simpleTable) align(data string, width int, direction int) string {
	strippedDataLen := len(table.stripAnsi(data))
	pad := width - strippedDataLen
	switch direction {
	case ALIGN_RIGHT:
		data = strings.Repeat(" ", pad) + data
	case ALIGN_CENTER:
		data = strings.Repeat(" ", pad/2) + data + strings.Repeat(" ", pad/2)
	default:
		data = data + strings.Repeat(" ", pad)
	}

	if len(data) != strippedDataLen {
		data += "\u001b[0m"
	}

	return data
}

func (table *simpleTable) renderCell(data string, width int, first bool, align int) string {
	// Trim data, if width is smaller
	if len(table.stripAnsi(data)) > width {
		data = data[:width-3-(table.padding*2)] + "..."
	}

	return table.align(strings.Repeat(" ", table.padding)+data+strings.Repeat(" ", table.padding), width, align)
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
	case _borderHeader:
		if table.style.inner.HEADER_IS_VISIBLE || table.style.inner.IS_VISIBLE {
			if table.style.outer.IS_VISIBLE {
				border = table.style.inner.HeaderLeft()
			}

			for idx, width := range rowWidths {
				border += strings.Repeat(table.style.inner.Header(), width)
				if idx < len(rowWidths)-1 {
					if table.style.inner.IS_VISIBLE {
						border += table.style.inner.HeaderMiddle()
					}
				} else {
					if table.style.outer.IS_VISIBLE {
						border += table.style.inner.HeaderRight()
					}
				}
			}
		}
	}
	return border
}

// Pivot data vertically to columns.
func (table *simpleTable) pivotData(data []string) [][]string {
	rowWidths := table.getRowWidths()
	cellBuff := make([][]string, len(data))
	maxrows := 0

	for cidx, cell := range data {
		if table.columnsTextWrap[cidx] {
			cellBuff[cidx] = textwrap.NewTextWrap().SetWidth(rowWidths[cidx] - (table.padding * 2)).Wrap(cell)
		} else {
			cellBuff[cidx] = []string{cell}
		}
		if len(cellBuff[cidx]) > maxrows {
			maxrows = len(cellBuff[cidx])
		}
	}

	pivoted := make([][]string, maxrows)
	for colIdx := 0; colIdx < maxrows; colIdx++ {
		pivotedRow := make([]string, len(rowWidths))
		for dataIdx, cellData := range cellBuff {
			if colIdx < len(cellData) {
				pivotedRow[dataIdx] = strings.TrimSpace(cellData[colIdx])
			} else {
				pivotedRow[dataIdx] = ""
			}
		}
		pivoted[colIdx] = pivotedRow
	}
	return pivoted
}

// Takes padded cells data and renders to the wrapped row
func (table *simpleTable) renderRowWrapped(cells []string) string {
	var rendered strings.Builder
	pivoted := table.pivotData(cells)
	for idx, innerRow := range pivoted {
		rendered.WriteString(table.renderRowSingle(innerRow))
		dlen := len(pivoted)
		if dlen > 1 && idx < dlen-1 {
			rendered.WriteString("\n")
		}
	}

	return rendered.String()
}

// Render row
func (table *simpleTable) renderRow(cells []string) string {
	var result string
	if table.wrapText {
		result = table.renderRowWrapped(cells)
	} else {
		result = table.renderRowSingle(cells)
	}

	return result
}

// Takes padded cells data and renders to the row with trimmed data
func (table *simpleTable) renderRowSingle(cells []string) string {
	rowWidths := table.getRowWidths()
	var row string
	for idx, cell := range cells {
		if idx < 1 {
			row += table.style.outer.VerticalLine()
		}
		row += table.renderCell(cell, rowWidths[idx], idx == 0, table.columnsAlign[idx])
		if idx < len(cells)-1 {
			row += table.style.inner.VerticalLine()
		} else {
			row += table.style.outer.VerticalLine()
		}
	}
	return row
}

// Renders table as a string
func (table *simpleTable) Render() string {
	table.setDataMaxWidth()
	render := make([]string, 0)

	if len(*table.Data().GetHeader()) > 0 {
		render = append(render, []string{
			table.renderBorder(_borderTop),
			table.renderRow(*table.Data().GetHeader()),
			table.renderBorder(_borderHeader),
		}...)
	}

	for _, row := range (*table.Data().GetData())[:len(*table.Data().GetData())-1] {
		render = append(render, []string{
			table.renderRow(row),
			table.renderBorder(_borderInner),
		}...)
	}
	render = append(render, []string{
		table.renderRow((*table.Data().GetData())[len(*table.Data().GetData())-1]),
		table.renderBorder(_borderBottom),
	}...)

	// Filter-out empty renders
	var rendered strings.Builder
	for _, chunk := range render {
		if len(chunk) > 0 {
			rendered.WriteString("\n" + chunk)
		}
	}
	return rendered.String()
}

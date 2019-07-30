package asciitable

import (
	"fmt"
	"github.com/isbm/textwrap"
	"regexp"
	"strconv"
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
	rowsData       *TableData
	rowsCount      uint64
	headerAlign    int
	cellAlign      int
	style          *borderStyle
	widthTable     int
	widthColumns   []int
	widthData      int
	padding        int
	wrapText       bool
	stripAnsiRegex *regexp.Regexp
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
	table.cellAlign = ALIGN_LEFT

	if style == nil {
		style = NewBorderStyle(-1, -1)
	}
	table.style = style

	table.widthTable, _ = GetTerminalSize()
	table.widthColumns = make([]int, 0)
	table.padding = 0
	table.wrapText = false
	table.stripAnsiRegex = regexp.MustCompile(_ansiRegex)
	table.SetDataMaxWidth()

	return table
}

// SetWrapText wraps text in all cells instead of trimming it to the max width.
func (table *simpleTable) SetWrapText(wrap bool) *simpleTable {
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
Set column width (chars)
*/
func (table *simpleTable) SetColWidth(column int, width int) *simpleTable {
	rowWidths := table.getRowWidths()
	if len(rowWidths) == 0 {
		panic("Attempt to set columns while no header or data has been set")
	}

	if len(table.widthColumns) == 0 {
		table.widthColumns = append(table.widthColumns, rowWidths...)
	}

	table.widthColumns[column] = width

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
func (table *simpleTable) SetDataMaxWidth() int {
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
func (table *simpleTable) getRowWidths() []int {
	widths := make([]int, len(*table.Data().GetHeader()))
	for idx := range widths {
		widths[idx] = 0
	}

	for idx, title := range *table.Data().GetHeader() {
		titleLength := len(table.stripAnsi(title))
		if titleLength > widths[idx] {
			widths[idx] = titleLength + table.padding*2
		}
	}

	for _, rData := range table.rowsData.data {
		for idx, data := range rData {
			dataLength := len(table.stripAnsi(data))
			if dataLength > widths[idx] {
				widths[idx] = dataLength + table.padding*2
			}
		}
	}

	// Override custom widths
	lastColumnFree := true
	if len(table.widthColumns) == len(widths) {
		lastColumnFree = widths[len(widths)-1] == table.widthColumns[len(widths)-1]
		copy(widths, table.widthColumns)
	}

	if lastColumnFree {
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
			widths[len(widths)-1] = widths[len(widths)-1] + table.widthTable - sum - offset
		} else {
			lastColWidth := table.widthTable - (sum - widths[len(widths)-1])
			if lastColWidth < 4 {
				lastColWidth = 4 + table.padding*2
			}
			widths[len(widths)-1] = lastColWidth
		}
	}

	return widths
}

func (table *simpleTable) renderCell(data string, width int, first bool) string {
	// Trim data, if width is smaller
	if len(table.stripAnsi(data)) > width {
		data = data[:width-3-(table.padding*2)] + "..."
	}

	// Set padding
	data = strings.Repeat(" ", table.padding) + data + strings.Repeat(" ", table.padding)

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

// Takes padded cells data and renders to the wrapped row
func (table *simpleTable) renderRowWrapped(cells []string) string {
	rowWidths := table.getRowWidths()
	cellBuff := make([][]string, len(cells))

	maxidx := 0
	maxrows := 0

	for cidx, cell := range cells {
		cellBuff[cidx] = textwrap.NewTextWrap().SetWidth(rowWidths[cidx] - (table.padding * 2)).Wrap(cell)
		if len(cellBuff[cidx]) > maxidx {
			maxrows = len(cellBuff[cidx])
			maxidx = cidx
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

	var rendered strings.Builder
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
		row += table.renderCell(cell, rowWidths[idx], idx == 0)
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
	render := make([]string, 0)

	if len(*table.Data().GetHeader()) > 0 {
		render = append(render, []string{
			table.renderBorder(_borderTop),
			table.renderRow(*table.Data().GetHeader()),
			table.renderBorder(_borderHeader),
		}...)
	}

	for _, row := range table.rowsData.data[:len(table.rowsData.data)-1] {
		render = append(render, []string{
			table.renderRow(row),
			table.renderBorder(_borderInner),
		}...)
	}
	render = append(render, []string{
		table.renderRow(table.rowsData.data[len(table.rowsData.data)-1]),
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

package asciitable

type borderOuter struct {
	VERTICAL_LINE   string
	HORISONTAL_LINE string
	LEFT_TOP        string
	LEFT_BOTTOM     string
	RIGHT_TOP       string
	RIGHT_BOTTOM    string
	IS_VISIBLE      bool
	style           int
}

type borderInner struct {
	VERTICAL_LINE     string
	HORISONTAL_LINE   string
	LEFT_MIDDLE       string
	CENTER_TOP        string
	CENTER_BOTTOM     string
	CENTER_MIDDLE     string
	RIGHT_MIDDLE      string
	HEADER_LEFT       string
	HEADER_MIDDLE     string
	HEADER_RIGHT      string
	HEADER            string
	HEADER_IS_VISIBLE bool
	IS_VISIBLE        bool
	style             int
}

type borderStyle struct {
	inner     borderInner
	outer     borderOuter
	widthFull bool
}

// Configuration of the table.
const (
	ALIGN_LEFT = iota
	ALIGN_CENTER
	ALIGN_RIGHT
	BORDER_SINGLE_THIN
	BORDER_SINGLE_THICK
	BORDER_DOUBLE
	BORDER_NONE
)

func NewBorderStyle(outer int, inner int) *borderStyle {
	style := new(borderStyle)

	style.outer = *new(borderOuter)
	style.outer.IS_VISIBLE = true
	style.outer.style = outer

	style.inner = *new(borderInner)
	style.inner.IS_VISIBLE = true
	style.inner.HEADER_IS_VISIBLE = true
	style.inner.style = inner
	style.widthFull = false

	style.initBorderStyle()

	return style
}

// Set table to be full width aligned to the terminal size.
// If this is set to False, then table will be calculated according to the
// data cells, if they are narrower than the terminal size.
func (style *borderStyle) SetTableWidthFull(full bool) *borderStyle {
	style.widthFull = full
	return style
}

func (style *borderStyle) initBorderStyle() *borderStyle {
	switch style.outer.style {
	case BORDER_SINGLE_THIN:
		style.outer.LEFT_TOP = "\u250c"
		style.outer.LEFT_BOTTOM = "\u2514"
		style.outer.RIGHT_TOP = "\u2510"
		style.outer.RIGHT_BOTTOM = "\u2518"
		style.outer.HORISONTAL_LINE = "\u2500"
		style.outer.VERTICAL_LINE = "\u2502"
		switch style.inner.style {
		case BORDER_SINGLE_THICK:
			style.inner.CENTER_TOP = "\u2530"
			style.inner.CENTER_BOTTOM = "\u2538"
			style.inner.CENTER_MIDDLE = "\u254b"
			style.inner.LEFT_MIDDLE = "\u251d"
			style.inner.RIGHT_MIDDLE = "\u2525"
			style.inner.HORISONTAL_LINE = "\u2501"
			style.inner.VERTICAL_LINE = "\u2503"
		case BORDER_DOUBLE:
			style.inner.CENTER_TOP = "\u2565"
			style.inner.CENTER_BOTTOM = "\u2568"
			style.inner.CENTER_MIDDLE = "\u256c"
			style.inner.LEFT_MIDDLE = "\u255e"
			style.inner.RIGHT_MIDDLE = "\u2561"
			style.inner.HORISONTAL_LINE = "\u2550"
			style.inner.VERTICAL_LINE = "\u2551"
		default: // Thin border
			style.inner.CENTER_TOP = "\u252c"
			style.inner.CENTER_BOTTOM = "\u2534"
			style.inner.CENTER_MIDDLE = "\u253c"
			style.inner.LEFT_MIDDLE = "\u251c"
			style.inner.RIGHT_MIDDLE = "\u2524"
			style.inner.HORISONTAL_LINE = "\u2500"
			style.inner.VERTICAL_LINE = "\u2502"
		}
	case BORDER_SINGLE_THICK:
		style.outer.LEFT_TOP = "\u250f"
		style.outer.LEFT_BOTTOM = "\u2517"
		style.outer.RIGHT_TOP = "\u2513"
		style.outer.RIGHT_BOTTOM = "\u251b"
		style.outer.HORISONTAL_LINE = "\u2501"
		style.outer.VERTICAL_LINE = "\u2503"

		switch style.inner.style {
		case BORDER_SINGLE_THICK:
			style.inner.CENTER_TOP = "\u2533"
			style.inner.CENTER_BOTTOM = "\u253b"
			style.inner.CENTER_MIDDLE = "\u254b"
			style.inner.LEFT_MIDDLE = "\u2523"
			style.inner.RIGHT_MIDDLE = "\u252b"
			style.inner.HORISONTAL_LINE = "\u2501"
			style.inner.VERTICAL_LINE = "\u2503"
		default: // Thin border, no double here
			style.inner.CENTER_TOP = "\u252f"
			style.inner.CENTER_BOTTOM = "\u2537"
			style.inner.CENTER_MIDDLE = "\u253c"
			style.inner.LEFT_MIDDLE = "\u2520"
			style.inner.RIGHT_MIDDLE = "\u2528"
			style.inner.HORISONTAL_LINE = "\u2500"
			style.inner.VERTICAL_LINE = "\u2502"
		}
	case BORDER_DOUBLE:
		style.outer.LEFT_TOP = "\u2554"
		style.outer.LEFT_BOTTOM = "\u255a"
		style.outer.RIGHT_TOP = "\u2557"
		style.outer.RIGHT_BOTTOM = "\u255d"
		style.outer.HORISONTAL_LINE = "\u2550"
		style.outer.VERTICAL_LINE = "\u2551"

		switch style.inner.style {
		case BORDER_DOUBLE:
			style.inner.CENTER_TOP = "\u2566"
			style.inner.CENTER_BOTTOM = "\u2569"
			style.inner.CENTER_MIDDLE = "\u256c"
			style.inner.LEFT_MIDDLE = "\u2560"
			style.inner.RIGHT_MIDDLE = "\u2563"
			style.inner.HORISONTAL_LINE = "\u2550"
			style.inner.VERTICAL_LINE = "\u2551"
		default: // Thin border, no thick here
			style.inner.CENTER_TOP = "\u2564"
			style.inner.CENTER_BOTTOM = "\u2567"
			style.inner.CENTER_MIDDLE = "\u253c"
			style.inner.LEFT_MIDDLE = "\u255f"
			style.inner.RIGHT_MIDDLE = "\u2562"
			style.inner.HORISONTAL_LINE = "\u2500"
			style.inner.VERTICAL_LINE = "\u2502"
		}
	default:
		// Just a plug ugly ascii style
		style.outer.LEFT_TOP = "+"
		style.outer.LEFT_BOTTOM = "+"
		style.outer.RIGHT_TOP = "+"
		style.outer.RIGHT_BOTTOM = "+"
		style.outer.HORISONTAL_LINE = "-"
		style.outer.VERTICAL_LINE = "|"
		style.inner.CENTER_TOP = "+"
		style.inner.CENTER_BOTTOM = "+"
		style.inner.CENTER_MIDDLE = "+"
		style.inner.LEFT_MIDDLE = "+"
		style.inner.RIGHT_MIDDLE = "+"
		style.inner.HORISONTAL_LINE = "-"
		style.inner.VERTICAL_LINE = "|"
	}

	return style
}

func (style *borderStyle) SetHeaderVisible(visibility bool) *borderStyle {
	style.inner.HEADER_IS_VISIBLE = visibility
	return style
}

// Set outer border visibility
func (style *borderStyle) SetBorderVisible(visibility bool) *borderStyle {
	style.outer.IS_VISIBLE = visibility
	if !style.outer.IS_VISIBLE {
		style.outer.LEFT_TOP, style.outer.LEFT_BOTTOM, style.outer.RIGHT_TOP, style.outer.RIGHT_BOTTOM,
			style.outer.HORISONTAL_LINE, style.outer.VERTICAL_LINE = "", "", "", "", "", ""
		style.inner.LEFT_MIDDLE, style.inner.RIGHT_MIDDLE, style.inner.CENTER_TOP,
			style.inner.CENTER_BOTTOM = "", "", "", ""
	} else {
		style.initBorderStyle()
	}

	return style
}

// Set table grid visibility
func (style *borderStyle) SetGridVisible(visibility bool) *borderStyle {
	style.inner.IS_VISIBLE = visibility
	if !style.inner.IS_VISIBLE {
		style.inner.CENTER_MIDDLE = ""
		style.inner.CENTER_BOTTOM = ""
		style.inner.CENTER_TOP = ""
		style.inner.VERTICAL_LINE = ""
		style.inner.HORISONTAL_LINE = ""
		style.inner.LEFT_MIDDLE = ""
		style.inner.RIGHT_MIDDLE = ""
	} else {
		style.initBorderStyle()
		style.SetBorderVisible(style.outer.IS_VISIBLE)
	}

	return style
}

// Set header style. This will draw a specified style line under the header.
func (style *borderStyle) SetHeaderStyle(header int) *borderStyle {
	switch header {
	case BORDER_SINGLE_THIN:
		if style.inner.IS_VISIBLE {
			switch style.inner.style {
			case BORDER_SINGLE_THIN:
				style.inner.HEADER = "\u2500"
				style.inner.HEADER_MIDDLE = "\u253c"
				switch style.outer.style {
				case BORDER_SINGLE_THIN:
					style.inner.HEADER_LEFT = "\u251c"
					style.inner.HEADER_RIGHT = "\u2524"
				case BORDER_SINGLE_THICK:
					style.inner.HEADER_LEFT = "\u2520"
					style.inner.HEADER_RIGHT = "\u2528"
				case BORDER_DOUBLE:
					style.inner.HEADER_LEFT = "\u255f"
					style.inner.HEADER_RIGHT = "\u2562"
				}
			}
		} else {
			style.inner.HEADER = "\u2500"
			style.inner.HEADER_LEFT = "\u2500"
			style.inner.HEADER_MIDDLE = "\u2500"
			style.inner.HEADER_RIGHT = "\u2500"
		}
	case BORDER_SINGLE_THICK:
		if style.inner.IS_VISIBLE {
			switch style.inner.style {
			case BORDER_SINGLE_THIN:
				style.inner.HEADER = "\u2501"
				style.inner.HEADER_MIDDLE = "\u253f"
				switch style.outer.style {
				case BORDER_SINGLE_THIN:
					style.inner.HEADER_LEFT = "\u251d"
					style.inner.HEADER_RIGHT = "\u2525"
				case BORDER_SINGLE_THICK:
					style.inner.HEADER_LEFT = "\u2523"
					style.inner.HEADER_RIGHT = "\u252b"
				case BORDER_DOUBLE:
					style.inner.HEADER_LEFT = "\u2560"
					style.inner.HEADER_RIGHT = "\u2563"
				}
			}
		} else {
			style.inner.HEADER = "\u2501"
			style.inner.HEADER_LEFT = "\u2501"
			style.inner.HEADER_MIDDLE = "\u2501"
			style.inner.HEADER_RIGHT = "\u2501"
		}
	case BORDER_DOUBLE:
		if style.inner.IS_VISIBLE {
			switch style.inner.style {
			case BORDER_SINGLE_THIN:
				style.inner.HEADER = "\u2550"
				style.inner.HEADER_MIDDLE = "\u256a"
				switch style.outer.style {
				case BORDER_SINGLE_THIN:
					style.inner.HEADER_LEFT = "\u255e"
					style.inner.HEADER_RIGHT = "\u2561"
				case BORDER_SINGLE_THICK:
					style.inner.HEADER_LEFT = "\u2503"
					style.inner.HEADER_MIDDLE = "\u2584"
					style.inner.HEADER_RIGHT = "\u2503"
					style.inner.HEADER = "\u2584"
				case BORDER_DOUBLE:
					style.inner.HEADER_LEFT = "\u2560"
					style.inner.HEADER_RIGHT = "\u2563"
				}
			}
		} else {
			style.inner.HEADER = "\u2550"
			style.inner.HEADER_LEFT = "\u2550"
			style.inner.HEADER_MIDDLE = "\u2550"
			style.inner.HEADER_RIGHT = "\u2550"
		}
	}
	return style
}

// Outer
func (border *borderOuter) LeftTop() string {
	return border.LEFT_TOP
}

func (border *borderOuter) RightTop() string {
	return border.RIGHT_TOP
}

func (border *borderOuter) LeftBottom() string {
	return border.LEFT_BOTTOM
}

func (border *borderOuter) RightBottom() string {
	return border.RIGHT_BOTTOM
}

func (border *borderOuter) HorisontalLine() string {
	return border.HORISONTAL_LINE
}

func (border *borderOuter) VerticalLine() string {
	return border.VERTICAL_LINE
}

// Inner
func (border *borderInner) LeftMiddle() string {
	return border.LEFT_MIDDLE
}

func (border *borderInner) CenterMiddle() string {
	return border.CENTER_MIDDLE
}
func (border *borderInner) RightMiddle() string {
	return border.RIGHT_MIDDLE

}
func (border *borderInner) CenterTop() string {
	return border.CENTER_TOP

}
func (border *borderInner) CenterBottom() string {
	return border.CENTER_BOTTOM

}
func (border *borderInner) HorisontalLine() string {
	return border.HORISONTAL_LINE

}
func (border *borderInner) VerticalLine() string {
	return border.VERTICAL_LINE
}

func (border *borderInner) Header() string {
	return border.HEADER
}

func (border *borderInner) HeaderLeft() string {
	return border.HEADER_LEFT
}

func (border *borderInner) HeaderMiddle() string {
	return border.HEADER_MIDDLE
}

func (border *borderInner) HeaderRight() string {
	return border.HEADER_RIGHT
}

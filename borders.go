package asciitable

type borderOuter struct {
	VERTICAL_LINE   string
	HORISONTAL_LINE string
	LEFT_TOP        string
	LEFT_BOTTOM     string
	RIGHT_TOP       string
	RIGHT_BOTTOM    string
	IS_VISIBLE      bool
}

type borderInner struct {
	VERTICAL_LINE   string
	HORISONTAL_LINE string
	LEFT_MIDDLE     string
	CENTER_TOP      string
	CENTER_BOTTOM   string
	CENTER_MIDDLE   string
	RIGHT_MIDDLE    string
	HEADER_LEFT     string
	HEADER_MIDDLE   string
	HEADER_RIGHT    string
	HEADER          string
	IS_VISIBLE      bool
}

type borderStyle struct {
	inner borderInner
	outer borderOuter
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

func NewBorderStyle(outer int, inner int, outerVisible bool, innerVisible bool) *borderStyle {
	style := new(borderStyle)

	style.outer = *new(borderOuter)
	style.outer.IS_VISIBLE = true

	style.inner = *new(borderInner)
	style.inner.IS_VISIBLE = true

	switch outer {
	case BORDER_SINGLE_THIN:
		style.outer.LEFT_TOP = "\u250c"
		style.outer.LEFT_BOTTOM = "\u2514"
		style.outer.RIGHT_TOP = "\u2510"
		style.outer.RIGHT_BOTTOM = "\u2518"
		style.outer.HORISONTAL_LINE = "\u2500"
		style.outer.VERTICAL_LINE = "\u2502"
		switch inner {
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

		switch inner {
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

		switch inner {
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

	style.outer.IS_VISIBLE = outerVisible
	style.inner.IS_VISIBLE = innerVisible

	if !style.outer.IS_VISIBLE {
		style.outer.LEFT_TOP, style.outer.LEFT_BOTTOM, style.outer.RIGHT_TOP, style.outer.RIGHT_BOTTOM,
			style.outer.HORISONTAL_LINE, style.outer.VERTICAL_LINE = "", "", "", "", "", ""
		style.inner.LEFT_MIDDLE, style.inner.RIGHT_MIDDLE, style.inner.CENTER_TOP,
			style.inner.CENTER_BOTTOM = "", "", "", ""
	}

	if !style.inner.IS_VISIBLE {
		style.inner.CENTER_MIDDLE = ""
		style.inner.CENTER_BOTTOM = ""
		style.inner.CENTER_TOP = ""
		style.inner.VERTICAL_LINE = ""
		style.inner.HORISONTAL_LINE = ""
		style.inner.LEFT_MIDDLE = ""
		style.inner.RIGHT_MIDDLE = ""
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

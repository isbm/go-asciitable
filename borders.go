package asciitable

type borderOuter struct {
	VERTICAL_LINE   string
	HORISONTAL_LINE string
	LEFT_TOP        string
	LEFT_BOTTOM     string
	RIGHT_TOP       string
	RIGHT_BOTTOM    string
}

type borderInner struct {
	VERTICAL_LINE   string
	HORISONTAL_LINE string
	LEFT_MIDDLE     string
	CENTER_TOP      string
	CENTER_BOTTOM   string
	CENTER_MIDDLE   string
	RIGHT_MIDDLE    string
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

func NewBorderStyle(outer int, inner int) *borderStyle {
	style := new(borderStyle)
	style.outer = *new(borderOuter)
	style.inner = *new(borderInner)

	switch outer {
	case BORDER_SINGLE_THIN:
		style.outer.LEFT_TOP = "\u250c"
		style.outer.LEFT_BOTTOM = "\u2514"
		style.outer.RIGHT_TOP = "\u2510"
		style.outer.RIGHT_BOTTOM = "\u2518"
		style.outer.HORISONTAL_LINE = "\u2500"
		style.outer.VERTICAL_LINE = "\u2502"

		switch inner {
		case BORDER_SINGLE_THIN:
			style.inner.CENTER_TOP = "\u252c"
			style.inner.CENTER_BOTTOM = "\u2534"
			style.inner.CENTER_MIDDLE = "\u253c"
			style.inner.LEFT_MIDDLE = "\u251c"
			style.inner.RIGHT_MIDDLE = "\u2524"
			style.inner.HORISONTAL_LINE = "\u2500"
			style.inner.VERTICAL_LINE = "\u2502"
		}
	default:
		// Plain single ascii
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

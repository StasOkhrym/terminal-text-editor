package editor

type Cursor struct {
	Row int
	Col int
}

func NewCursor() *Cursor {
	// The cursor is initially positioned with padding from header
	return &Cursor{
		Row: 0,
		Col: 2,
	}
}

func (c *Cursor) MoveTo(row int, col int) {
	c.Row = row
	c.Col = col
}

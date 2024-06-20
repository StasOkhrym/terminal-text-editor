package editor

func (e *Editor) MoveCursorTo(dx, dy int) {
	headerRows := 2
	footerRows := 2
	totalRows := e.tty.WindowSize().Rows
	totalCols := e.tty.WindowSize().Cols

	newRow := e.cursor.Row + dy
	newCol := e.cursor.Col + dx

	if newRow <= headerRows || newRow >= totalRows-footerRows || newCol < 0 || newCol >= totalCols {
		return
	}

	e.tty.MoveCursorTo(newCol, newRow)
}

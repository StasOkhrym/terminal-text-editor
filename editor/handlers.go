package editor

import (
	"fmt"
)

func (e *Editor) MoveCursorTo(dx, dy int) {
	headerRows := 2
	footerRows := 2
	totalRows := e.tty.WindowSize().Rows
	totalCols := e.tty.WindowSize().Cols

	newRow := e.cursor.Row + dy
	newCol := e.cursor.Col + dx

	// Check if the new cursor position is within the total rows and columns, excluding the header and footer rows
	if newRow < headerRows || newRow >= totalRows-footerRows || newCol < 0 || newCol >= totalCols {
		e.buffer.AppendToBuffer([]byte(fmt.Sprintf("Cursor out of bounds: %d, %d\n", newRow, newCol)))
		// return
	}

	e.cursor.Row = newRow
	e.cursor.Col = newCol

	e.tty.MoveCursorTo(newCol, newRow)
}

func (e *Editor) SaveFile() error {
	_, err := e.file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = e.file.Write(e.buffer.Output())
	if err != nil {
		return err
	}

	return nil
}

func (e *Editor) Close() error {
	if e.tty == nil && e.file == nil {
		return nil
	}

	if err := e.tty.Close(); err != nil {
		return err
	}
	e.tty = nil

	if err := e.file.Close(); err != nil {
		return err
	}
	e.file = nil

	return nil
}

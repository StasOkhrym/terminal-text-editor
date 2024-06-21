package editor

import (
	"os"
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
		// e.buffer.AppendToBuffer([]byte(fmt.Sprintf("Cursor out of bounds: %d, %d\n", newRow, newCol)))
		return
	}

	e.cursor.Row = newRow
	e.cursor.Col = newCol

	e.tty.MoveCursorTo(newCol, newRow)
}

func (e *Editor) SaveFile() error {
	return os.WriteFile(e.file.Name(), []byte(e.buffer.String()), 0644)

}

func (e *Editor) Close() error {
	if e.closed {
		return nil
	}
	close(e.receiver)
	e.closed = true

	if e.tty == nil && e.file == nil {
		return nil
	}

	if err := e.tty.Close(); err != nil {
		return err
	}

	if err := e.file.Close(); err != nil {
		return err
	}

	return nil
}

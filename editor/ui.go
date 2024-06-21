package editor

import (
	"fmt"
)

type ScreenBuffer struct {
	Rows [][]byte
}

func NewScreenBuffer(rows, cols int) *ScreenBuffer {
	buffer := make([][]byte, rows)
	for i := range buffer {
		buffer[i] = make([]byte, cols)
	}
	return &ScreenBuffer{Rows: buffer}
}

func (sb *ScreenBuffer) Set(row, col int, b byte) {
	sb.Rows[row][col] = b
}

func (sb *ScreenBuffer) Get(row, col int) byte {
	return sb.Rows[row][col]
}

func (e *Editor) RenderUI(filePath string) {
	totalRows := e.tty.WindowSize().Rows
	totalCols := e.tty.WindowSize().Cols

	headerRows := 2

	e.tty.Output().WriteString("\033[2J")

	e.tty.MoveCursorTo(0, 0)
	headerText := fmt.Sprintf("File: %s", filePath)
	headerPadding := (totalCols - len(headerText)) / 2
	paddedHeaderText := fmt.Sprintf("%*s%s%*s", headerPadding, "", headerText, headerPadding, "")
	e.tty.Output().WriteString(fmt.Sprintf("\033[44m\033[37m%*s\033[0m", totalCols, paddedHeaderText))

	e.tty.MoveCursorTo(headerRows, 0)
	e.tty.Output().WriteString(string(e.buffer.String()))

	e.tty.MoveCursorTo(totalRows-1, 0)
	footerText := fmt.Sprintf(
		"CTRL + S: Save      CTRL + X: Exit   Cursor Row: %d   Cursor Col: %d",
		e.cursor.Row,
		e.cursor.Col,
	)
	padding := (totalCols - len(footerText)) / 2
	paddedFooterText := fmt.Sprintf("%*s%s%*s", padding, "", footerText, padding, "")
	e.tty.Output().WriteString(fmt.Sprintf("\033[42m\033[37m%*s\033[0m", totalCols, paddedFooterText))

	e.tty.MoveCursorTo(e.cursor.Row, e.cursor.Col)
}

func (e *Editor) RenderSaved(filePath string) {
	totalRows := e.tty.WindowSize().Rows
	totalCols := e.tty.WindowSize().Cols

	headerRows := 2

	e.tty.Output().WriteString("\033[2J")

	e.tty.MoveCursorTo(0, 0)
	headerText := fmt.Sprintf("File: %s", filePath)
	headerPadding := (totalCols - len(headerText)) / 2
	paddedHeaderText := fmt.Sprintf("%*s%s%*s", headerPadding, "", headerText, headerPadding, "")
	e.tty.Output().WriteString(fmt.Sprintf("\033[44m\033[37m%*s\033[0m", totalCols, paddedHeaderText))

	e.tty.MoveCursorTo(headerRows, 0)
	e.tty.Output().WriteString(string(e.buffer.String()))

	e.tty.MoveCursorTo(totalRows-1, 0)
	footerText := "File Saved! Press any key to continue."
	padding := (totalCols - len(footerText)) / 2
	paddedFooterText := fmt.Sprintf("%*s%s%*s", padding, "", footerText, padding, "")
	e.tty.Output().WriteString(fmt.Sprintf("\033[42m\033[37m%*s\033[0m", totalCols, paddedFooterText))

	e.tty.MoveCursorTo(e.cursor.Row, e.cursor.Col)
}

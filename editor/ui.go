package editor

import (
	"fmt"
)

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
	e.tty.Output().WriteString(string(e.buffer.Buffer()))

	e.tty.MoveCursorTo(totalRows-1, 0)
	footerText := "CTRL + S: Save      CTRL + X: Exit"
	padding := (totalCols - len(footerText)) / 2
	paddedFooterText := fmt.Sprintf("%*s%s%*s", padding, "", footerText, padding, "")
	e.tty.Output().WriteString(fmt.Sprintf("\033[42m\033[37m%*s\033[0m", totalCols, paddedFooterText))

	e.tty.MoveCursorTo(e.cursor.Col, e.cursor.Row)
}

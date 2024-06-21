package editor

import (
	"os"
	"os/signal"
	"syscall"

	"text_editor/buffer"
	"text_editor/teletype"
)

type Editor struct {
	buffer *buffer.Matrix
	file   *os.File
	tty    *teletype.TTY

	cursor *Cursor
}

func NewEditor(filePath string) (*Editor, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	tty, err := teletype.Open()
	if err != nil {
		file.Close()
		return nil, err
	}
	defer tty.Close()

	rows, cols := tty.WindowSize().Rows, tty.WindowSize().Cols
	m := buffer.NewMatrix(rows, cols, content)

	return &Editor{
		buffer: m,
		file:   file,
		tty:    tty,
		cursor: NewCursor(),
	}, nil
}

func (e *Editor) Run() error {
	defer e.Close()

	e.tty.EnableAlternateScreenBuffer()
	defer e.tty.DisableAlternateScreenBuffer()

	e.RenderUI(e.file.Name(), false)

	winResize := make(chan os.Signal, 1)
	signal.Notify(winResize, syscall.SIGWINCH)

	go func() {
		for range winResize {
			e.tty.UpdateWindowSize()
			e.RenderUI(e.file.Name(), false)
		}
	}()

	// Main event loop
	for {
		bytes, err := e.tty.ReadKey()
		if err != nil {
			os.Exit(1)
		}

		switch string(bytes) {
		case teletype.KeyArrowUp:
			e.MoveCursorTo(0, -1)
		case teletype.KeyArrowDown:
			e.MoveCursorTo(0, 1)
		case teletype.KeyArrowRight:
			e.MoveCursorTo(1, 0)
		case teletype.KeyArrowLeft:
			e.MoveCursorTo(-1, 0)
		case teletype.KeyCtrlS:
			e.RenderUI(e.file.Name(), true)
			e.SaveFile()
			for {
				upBytes, err := e.tty.ReadKey()
				if err != nil {
					os.Exit(1)
				}

				if len(upBytes) != 0 {
					break
				}
			}

		case teletype.KeyCtrlX:
			return e.Close()
		case teletype.KeyBackspace:
			e.buffer.Pop(e.cursor.Row-2, e.cursor.Col)
			e.MoveCursorTo(-1, 0)
		default:
			e.buffer.Set(e.cursor.Row-2, e.cursor.Col, bytes)
			e.MoveCursorTo(1, 0)
		}

		// Refresh UI after handling the key
		e.RenderUI(e.file.Name(), false)
	}
}

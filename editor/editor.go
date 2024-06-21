package editor

import (
	"os"
	"os/signal"
	"syscall"

	"text_editor/buffer"
	"text_editor/teletype"
)

type Editor struct {
	buffer *buffer.GapBuffer
	file   *os.File
	tty    *teletype.TTY

	cursor *Cursor
}

func NewEditor(filePath string) (*Editor, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	gb := buffer.NewGapBuffer(1024)
	_, err = file.Read(gb.Buffer())
	if err != nil {
		return nil, err
	}

	tty, err := teletype.Open()
	if err != nil {
		file.Close()
		return nil, err
	}
	defer tty.Close()

	return &Editor{
		buffer: gb,
		file:   file,
		tty:    tty,
		cursor: NewCursor(),
	}, nil
}

func (e *Editor) Run() error {
	defer e.Close()

	e.tty.EnableAlternateScreenBuffer()
	defer e.tty.DisableAlternateScreenBuffer()

	e.RenderUI(e.file.Name())

	winResize := make(chan os.Signal, 1)
	signal.Notify(winResize, syscall.SIGWINCH)

	go func() {
		for range winResize {
			e.RenderUI(e.file.Name())
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
			e.SaveFile()
			e.buffer.AppendToBuffer([]byte("Saving file...\n"))
		case teletype.KeyCtrlX:
			return e.Close()
		}

		// Refresh UI after handling the key
		e.RenderUI(e.file.Name())
	}
}

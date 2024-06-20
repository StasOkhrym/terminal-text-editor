package editor

import (
	"os"

	"text_editor/buffer"
	"text_editor/teletype"
)

type Editor struct {
	buffer *buffer.GapBuffer
	file   *os.File
	tty    *teletype.TTY
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
	}, nil
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

func (e *Editor) Run() error {
	defer e.Close()

	e.tty.EnableAlternateScreenBuffer()
	defer e.tty.DisableAlternateScreenBuffer()

	e.tty.Output().WriteString("Text editor. Press Ctrl+X to exit.\n\n")
	e.tty.Output().WriteString(e.buffer.Text())

	for {
		bytes, err := e.tty.ReadKey()
		e.tty.Output().WriteString(string(bytes))
		if err != nil {
			os.Exit(1)
		}

		switch string(bytes) {
		case teletype.KeyArrowUp:
			e.tty.UpdateCursorPosition(0, -1)
		case teletype.KeyArrowDown:
			e.tty.UpdateCursorPosition(0, 1)
		case teletype.KeyArrowRight:
			e.tty.UpdateCursorPosition(1, 0)
		case teletype.KeyArrowLeft:
			e.tty.UpdateCursorPosition(-1, 0)
		case teletype.KeyCtrlS:
			e.tty.Output().WriteString("pasting")
		// case teletype.KeyDelete:
		// 	e.tty.Delete()
		// case teletype.KeyEnter:
		// 	e.tty.Insert('\n')
		case teletype.KeyCtrlX:
			return e.Close()
		}
	}

}

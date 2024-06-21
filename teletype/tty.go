package teletype

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

type WindowSize struct {
	Rows int
	Cols int
}

// TTY base representation for Teletypewriter
type TTY struct {
	in         *os.File
	bin        *bufio.Reader
	out        *os.File
	termios    unix.Termios
	windowSize *WindowSize
}

func Open() (*TTY, error) {
	return open("/dev/tty")
}

func OpenDevice(path string) (*TTY, error) {
	return open(path)
}

func open(path string) (*TTY, error) {
	in, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	bin := bufio.NewReader(in)

	out, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}

	termios, err := unix.IoctlGetTermios(int(in.Fd()), unix.TIOCGETA)
	if err != nil {
		return nil, err
	}

	// Applying bitmask for Termios settings
	termios.Iflag &^= unix.ISTRIP | unix.INLCR | unix.ICRNL | unix.IGNCR | unix.IXON
	termios.Lflag &^= unix.ECHO | unix.ICANON
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	// Have no idea why it is still processing backspace
	termios.Cc[unix.VERASE] = 0
	termios.Cc[unix.VWERASE] = 0
	termios.Cc[unix.VINTR] = 0

	err = unix.IoctlSetTermios(int(in.Fd()), unix.TIOCSETA, termios)
	if err != nil {
		return nil, err
	}

	tty := &TTY{
		in:      in,
		bin:     bin,
		out:     out,
		termios: *termios,
		windowSize: &WindowSize{
			Rows: 0,
			Cols: 0,
		},
	}

	tty.UpdateWindowSize()

	return tty, nil
}

func (tty *TTY) UpdateWindowSize() {
	ws, err := unix.IoctlGetWinsize(int(tty.in.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}

	tty.windowSize.Rows = int(ws.Row)
	tty.windowSize.Cols = int(ws.Col)
}

func (tty *TTY) Close() error {
	return unix.IoctlSetTermios(int(tty.in.Fd()), unix.TIOCSETA, &tty.termios)
}

func (tty *TTY) Input() *os.File {
	return tty.in
}

func (tty *TTY) InputReader() *bufio.Reader {
	return tty.bin
}

func (tty *TTY) Output() *os.File {
	return tty.out
}

func (tty *TTY) WindowSize() *WindowSize {
	return tty.windowSize
}

func (tty *TTY) EnableAlternateScreenBuffer() {
	tty.out.WriteString("\033[?1049h")
}

func (tty *TTY) DisableAlternateScreenBuffer() {
	tty.out.WriteString("\033[?1049l")
}

func (tty *TTY) ClearScreen() {
	tty.out.WriteString("\033[2J\033[H")
}

func (tty *TTY) Cleanup() error {
	tty.DisableAlternateScreenBuffer()
	// tty.ClearScreen()
	return tty.Close()
}

func (tty *TTY) MoveCursorTo(row int, col int) {
	if row < 0 || col < 0 || col >= tty.windowSize.Cols || row >= tty.windowSize.Rows {
		return
	}
	escapeCode := fmt.Sprintf("\033[%d;%dH", row+1, col+1)
	tty.out.WriteString(escapeCode)
}

// ReadKey reads user input including escape sequences
func (tty *TTY) ReadKey() ([]byte, error) {
	b, err := tty.bin.ReadByte()
	if err != nil {
		return nil, err
	}

	if b == '\x1b' {
		seq := make([]byte, 0, 3)
		seq = append(seq, b)

		for i := 0; i < 2; i++ {
			b, err = tty.bin.ReadByte()
			if err != nil {
				return nil, err
			}
			seq = append(seq, b)
		}

		return seq, nil
	}

	return []byte{b}, nil
}

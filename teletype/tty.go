package teletype

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	sig        chan os.Signal
	windowSize *WindowSize

	closed bool
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

	err = unix.IoctlSetTermios(int(in.Fd()), unix.TIOCSETA, termios)
	if err != nil {
		return nil, err
	}

	sig := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGINT)
	signal.Notify(sig, syscall.SIGWINCH)

	tty := &TTY{
		in:      in,
		bin:     bin,
		out:     out,
		termios: *termios,
		sig:     sig,
		windowSize: &WindowSize{
			Rows: 0,
			Cols: 0,
		},

		closed: false,
	}

	go func() {
		for range sig {
			tty.UpdateWindowSize()
		}
	}()

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
	if tty.closed {
		return nil
	}
	signal.Stop(tty.sig)
	if tty.sig != nil {
		close(tty.sig)
		tty.sig = nil
	}
	tty.closed = true
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

func (tty *TTY) Signal() chan os.Signal {
	return tty.sig
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
	if row < 0 || col < 0 || row >= tty.windowSize.Cols || col >= tty.windowSize.Rows {
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

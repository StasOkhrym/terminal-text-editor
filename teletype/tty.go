package teletype

import (
	"bufio"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

// TTY base representation for Teletypewriter
type TTY struct {
	in      *os.File
	bin     *bufio.Reader
	out     *os.File
	termios unix.Termios
	sig     chan os.Signal
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
	termios.Iflag &^= unix.ISTRIP | unix.INLCR | unix.ICRNL | unix.IGNCR | unix.IXOFF
	termios.Lflag &^= unix.ECHO | unix.ICANON
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(int(in.Fd()), unix.TIOCSETA, termios)
	if err != nil {
		return nil, err
	}

	sig := make(chan os.Signal, 1)

	return &TTY{
		in:      in,
		bin:     bin,
		out:     out,
		termios: *termios,
		sig:     sig,
	}, nil
}

func (tty *TTY) Close() error {
	signal.Stop(tty.sig)
	close(tty.sig)
	return unix.IoctlSetTermios(int(tty.in.Fd()), unix.TIOCSETA, &tty.termios)
}

func (tty *TTY) Read() (rune, error) {
	r, _, err := tty.bin.ReadRune()
	return r, err
}

func (tty *TTY) Input() *os.File {
	return tty.in
}

func (tty *TTY) Output() *os.File {
	return tty.out
}

func (tty *TTY) Signal() chan os.Signal {
	return tty.sig
}

package teletype

const (
	// Cursor control characters
	KeyArrowUp    = "\x1b[A"
	KeyArrowDown  = "\x1b[B"
	KeyArrowRight = "\x1b[C"
	KeyArrowLeft  = "\x1b[D"

	// Text editing characters
	KeyBackspace = "\x08"
	KeyDelete    = "\x7f"
	KeyEnter     = "\r"

	// Control characters
	KeyCtrlS = "\x13"
	KeyCtrlZ = "\x1a"

	// Exit
	KeyCtrlX = "\x18"
)

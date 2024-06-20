package buffer

type GapBuffer struct {
	buffer   []byte
	gapStart int
	gapEnd   int
}

func NewGapBuffer(initialSize int) *GapBuffer {
	return &GapBuffer{
		buffer:   make([]byte, initialSize),
		gapStart: 0,
		gapEnd:   initialSize,
	}
}

func (gb *GapBuffer) Insert(r rune) {
	if gb.gapStart == gb.gapEnd {
		gb.grow()
	}
	gb.buffer[gb.gapStart] = byte(r)
	gb.gapStart++
}

func (gb *GapBuffer) Delete() {
	if gb.gapStart > 0 {
		gb.gapStart--
	}
}

func (gb *GapBuffer) grow() {
	newSize := len(gb.buffer) * 2
	newBuffer := make([]byte, newSize)
	copy(newBuffer, gb.buffer[:gb.gapStart])
	copy(newBuffer[gb.gapStart+newSize/2:], gb.buffer[gb.gapEnd:])
	gb.buffer = newBuffer
	gb.gapEnd += newSize / 2
}

func (gb *GapBuffer) Text() string {
	return string(gb.buffer[:gb.gapStart]) + string(gb.buffer[gb.gapEnd:])
}

func (gb *GapBuffer) MoveGapTo(index int) {
	if index < gb.gapStart {
		copy(gb.buffer[index+gb.gapEnd-gb.gapStart:], gb.buffer[index:gb.gapStart])
		gb.gapEnd += gb.gapStart - index
		gb.gapStart = index
	} else if index > gb.gapStart {
		copy(gb.buffer[gb.gapStart:], gb.buffer[gb.gapEnd:index])
		gb.gapStart += index - gb.gapEnd
		gb.gapEnd = index
	}
}

func (gb *GapBuffer) GapStart() int {
	return gb.gapStart
}

func (gb *GapBuffer) GapEnd() int {
	return gb.gapEnd
}

func (gb *GapBuffer) Buffer() []byte {
	return gb.buffer
}

func (gb *GapBuffer) Size() int {
	return len(gb.buffer)
}

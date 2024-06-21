package buffer

import (
	"bytes"
	"strings"
)

type Matrix struct {
	Rows [][]byte
}

func NewMatrix(rows, cols int, content []byte) *Matrix {
	matrix := make([][]byte, rows)
	lines := bytes.Split(content, []byte("\n"))
	for i := range matrix {
		if i < len(lines) {
			line := lines[i]
			if len(line) > cols {
				line = line[:cols]
			} else if len(line) < cols {
				line = append(line, bytes.Repeat([]byte(" "), cols-len(line))...)
			}
			matrix[i] = line
		} else {
			matrix[i] = make([]byte, cols)
		}
	}
	return &Matrix{Rows: matrix}
}

func (m *Matrix) Set(row, col int, b []byte) {
	for len(m.Rows) <= row {
		m.Rows = append(m.Rows, make([]byte, 0))
	}

	if len(m.Rows[row]) >= col+len(b) {
		copy(m.Rows[row][col:], b)
	} else {
		newRow := make([]byte, col+len(b))
		copy(newRow, m.Rows[row][:col])
		copy(newRow[col:], b)
		m.Rows[row] = newRow
	}
}

func (m *Matrix) Get(row, col int) byte {
	return m.Rows[row][col]
}

func (m *Matrix) Pop(row, col int) {
	if col > 0 && len(m.Rows[row]) > 0 {
		m.Rows[row] = append(m.Rows[row][:col-1], m.Rows[row][col:]...)
	}
}

func (m *Matrix) String() string {
	var lines []string
	for _, row := range m.Rows {
		row = bytes.TrimRight(row, " ")
		row = bytes.Trim(row, "\x00")
		lines = append(lines, string(row))
	}

	// Remove trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n")
}

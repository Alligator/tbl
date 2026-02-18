// Package tbl implements a simple way to print tables
package tbl

import (
	"fmt"
	"regexp"
	"strings"
)

type Table struct {
	rows        [][]strings.Builder
	colNames    []string
	currentCols int
	maxCols     int
}

func NewTable() *Table {
	return &Table{
		rows:     make([][]strings.Builder, 0),
		colNames: make([]string, 32),
	}
}

// NewRow adds a row to the end of table
func (t *Table) NewRow() {
	cols := make([]strings.Builder, 0)
	t.rows = append(t.rows, cols)
	t.currentCols = 0
}

// NewCol adds a column with the given name to the current row
func (t *Table) NewCol(name string) {
	curRow := &t.rows[len(t.rows)-1]
	*curRow = append(*curRow, strings.Builder{})
	t.colNames[t.currentCols] = name
	t.currentCols++
	if t.currentCols > t.maxCols {
		t.maxCols = t.currentCols
	}
}

// Print s to the current column
func (t *Table) Print(s string) {
	curRow := &t.rows[len(t.rows)-1]
	(*curRow)[len(*curRow)-1].WriteString(s)
}

// Printf is shorthand for
//
//	tbl.Print(fmt.Sprintf(format, a))
func (t *Table) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	t.Print(s)
}

var escapeRegex = regexp.MustCompile(`\x1b((\[\d+m)|(\]8;;[^\x1b]+\x1b\\)|(\]8;;\x1b\\))`)

func lenWithoutEscapes(s string) int {
	rs := escapeRegex.ReplaceAllString(s, "")
	return len(rs)
}

// String returns the formatted table
func (t *Table) String() string {
	colWidths := make([]int, t.maxCols)
	for _, row := range t.rows {
		for i, col := range row {
			cl := lenWithoutEscapes(col.String())
			if colWidths[i] < cl {
				colWidths[i] = cl
			}
		}
	}

	for i := range colWidths {
		hl := lenWithoutEscapes(t.colNames[i])
		if hl > colWidths[i] {
			colWidths[i] = hl
		}
	}

	var sb strings.Builder

	// header
	sb.WriteString("| ")
	for i := range colWidths {
		sb.WriteString(fmt.Sprintf("%-*s |", colWidths[i], t.colNames[i]))
		if i < len(colWidths)-1 {
			sb.WriteRune(' ')
		}
	}
	sb.WriteString("\n")

	sb.WriteString("| ")
	for i := range colWidths {
		sb.WriteString(strings.Repeat("-", colWidths[i]))
		sb.WriteString(" |")
		if i < len(colWidths)-1 {
			sb.WriteRune(' ')
		}
	}
	sb.WriteString("\n")

	// body
	for _, row := range t.rows {
		sb.WriteString("| ")
		for i, col := range row {
			s := col.String()
			lwe := lenWithoutEscapes(s)
			w := colWidths[i] + (len(s) - lwe)
			sb.WriteString(fmt.Sprintf("%-*s |", w, s))
			if i < len(row)-1 {
				sb.WriteRune(' ')
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

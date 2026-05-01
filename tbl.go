// Package tbl implements a simple way to print tables
package tbl

import (
	"fmt"
	"regexp"
	"strings"
)

type TableStyle uint8

const (
	StyleDefault TableStyle = iota
	StyleMinimal
)

type Table struct {
	Style       TableStyle
	rows        [][]strings.Builder
	colNames    []string
	currentCols int
	maxCols     int
}

type styleProps struct {
	colStart        string
	colSep          string
	colEnd          string
	headerSep       string
	uppercaseHeader bool
}

var styles map[TableStyle]styleProps = map[TableStyle]styleProps{
	StyleDefault: {
		colStart:  "| ",
		colSep:    " | ",
		colEnd:    " |",
		headerSep: "-",
	},
	StyleMinimal: {
		colStart:        "",
		colSep:          "  ",
		colEnd:          "",
		headerSep:       "",
		uppercaseHeader: true,
	},
}

func NewTable() *Table {
	return &Table{
		Style:    StyleDefault,
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

func (t *Table) getRow() *[]strings.Builder {
	if len(t.rows) == 0 {
		t.NewRow()
	}
	return &t.rows[len(t.rows)-1]
}

// NewCol adds a column with the given name to the current row
func (t *Table) NewCol(name string) {
	curRow := t.getRow()
	*curRow = append(*curRow, strings.Builder{})
	t.colNames[t.currentCols] = name
	t.currentCols++
	if t.currentCols > t.maxCols {
		t.maxCols = t.currentCols
	}
}

// Print s to the current column
func (t *Table) Print(s string) {
	curRow := t.getRow()
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

	styl := styles[t.Style]

	// header
	sb.WriteString(styl.colStart)
	for i := range colWidths {
		header := t.colNames[i]
		if styl.uppercaseHeader {
			header = strings.ToUpper(header)
		}

		if i < len(colWidths)-1 {
			sb.WriteString(fmt.Sprintf("%-*s%s", colWidths[i], header, styl.colSep))
		} else if len(styl.colEnd) > 0 {
			sb.WriteString(fmt.Sprintf("%-*s%s", colWidths[i], header, styl.colEnd))
		} else {
			sb.WriteString(header)
		}
	}
	sb.WriteString("\n")

	if len(styl.headerSep) > 0 {
		sb.WriteString(styl.colStart)
		for i := range colWidths {
			sb.WriteString(strings.Repeat(styl.headerSep, colWidths[i]))
			if i < len(colWidths)-1 {
				sb.WriteString(styl.colSep)
			}
		}
		sb.WriteString(styl.colEnd)
		sb.WriteString("\n")
	}

	// body
	for _, row := range t.rows {
		sb.WriteString(styl.colStart)
		for i, col := range row {
			s := col.String()
			lwe := lenWithoutEscapes(s)
			w := colWidths[i] + (len(s) - lwe)
			if i < len(colWidths)-1 {
				sb.WriteString(fmt.Sprintf("%-*s", w, s))
				sb.WriteString(styl.colSep)
			} else if len(styl.colEnd) > 0 {
				sb.WriteString(fmt.Sprintf("%-*s", w, s))
				sb.WriteString(styl.colEnd)
			} else {
				sb.WriteString(s)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

package table

import (
	"fmt"
	"imageserver/internal/model"
	"strings"
)

type Lengths struct {
	FileName int
	Created  int
	Updated  int
}

func NewLengths(lf model.ListFile) Lengths {
	sl := MeasureStrings(lf)
	return Lengths{
		FileName: sl.FileName,
		Created:  sl.Created,
		Updated:  sl.Updated,
	}
}

func (f *Form) AddToForm(s string) {
	f.form = append(f.form, s)
}

type Form struct {
	form []string
}

// MakeTable converts model.ListFile to a table of type string.
func MakeTable(lf model.ListFile) string {
	length := NewLengths(lf)
	f := &Form{}
	f.AddToForm(fmt.Sprintf("%c%s%c%s%c%s%c", table.LeftTop, table.RepeatLine(length.FileName), table.CenterTop, table.RepeatLine(length.Created), table.CenterTop, table.RepeatLine(length.Updated), table.RightTop))
	f.AddToForm(fmt.Sprintf("%c%s%c%s%c%s%c", table.Vertical, Fit(table.FileName, length.FileName), table.Vertical, Fit(table.Created, length.Created), table.Vertical, Fit(table.Updated, length.Updated), table.Vertical))
	f.AddToForm(fmt.Sprintf("%c%s%c%s%c%s%c", table.LeftMiddle, table.RepeatLine(length.FileName), table.CenterMiddle, table.RepeatLine(length.Created), table.CenterMiddle, table.RepeatLine(length.Updated), table.RightMiddle))
	for _, v := range lf {
		f.AddToForm(fmt.Sprintf("%c%s%c%s%c%s%c", table.Vertical, Fit(v.FileName, length.FileName), table.Vertical, Fit(v.Created, length.Created), table.Vertical, Fit(v.Updated, length.Updated), table.Vertical))
	}
	f.AddToForm(fmt.Sprintf("%c%s%c%s%c%s%c", table.LeftBottom, table.RepeatLine(length.FileName), table.CenterBottom, table.RepeatLine(length.Created), table.CenterBottom, table.RepeatLine(length.Updated), table.RightBottom))
	result := strings.Join(f.form, "\n")
	return result
}

// MeasureStrings measures the maximum lengths of the table header names to change the table width.
func MeasureStrings(lf model.ListFile) Lengths {
	var fileName, created, updated int
	air := 2
	for _, v := range lf {
		if len(v.FileName) > fileName {
			fileName = len(v.FileName)
		}
		if len(v.Created) > created {
			created = len(v.Created)
		}
		if len(v.Updated) > updated {
			updated = len(v.Updated)
		}
	}
	if fileName < len(table.FileName) {
		fileName = len(table.FileName)
	}
	if created < len(table.Created) {
		created = len(table.Created)
	}
	if updated < len(table.Updated) {
		updated = len(table.Updated)
	}
	return Lengths{fileName + air, created + air, updated + air}
}

// Fit adds spaces in the table header name or in the names of files that are in the first column.
func Fit(s string, n int) string {
	for len(s) < n {
		s = fmt.Sprintf("%s%c", s, table.WhiteSpace)
		if len(s) == n {
			break
		}
		s = fmt.Sprintf("%c%s", table.WhiteSpace, s)
		if len(s) == n {
			break
		}
	}
	return s
}

type Table struct {
	WhiteSpace   rune
	LeftTop      rune
	CenterTop    rune
	RightTop     rune
	Vertical     rune
	Horizontal   rune
	LeftMiddle   rune
	CenterMiddle rune
	RightMiddle  rune
	LeftBottom   rune
	CenterBottom rune
	RightBottom  rune
	FileName     string
	Created      string
	Updated      string
}

var table = NewTable()

// RepeatLine repeats the "═" symbol to draw the horizontal lines of the table.
func (t Table) RepeatLine(n int) string {
	return strings.Repeat(string(t.Horizontal), n)
}

func NewTable() Table {
	return Table{
		WhiteSpace:   '\u0020', // " "
		LeftTop:      '\u2554', // "╔"
		CenterTop:    '\u2566', // "╦"
		RightTop:     '\u2557', // "╗"
		Vertical:     '\u2551', // "║"
		Horizontal:   '\u2550', // "═"
		LeftMiddle:   '\u2560', // "╠"
		CenterMiddle: '\u256c', // "╬"
		RightMiddle:  '\u2563', // "╣"
		LeftBottom:   '\u255a', // "╚"
		CenterBottom: '\u2569', // "╩"
		RightBottom:  '\u255d', // "╝"
		FileName:     "File name",
		Created:      "Created",
		Updated:      "Last updated",
	}
}

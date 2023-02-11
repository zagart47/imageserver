package table

import (
	"fmt"
	"imageserver/file"
	"strings"
)

type Lengths struct {
	FileNameLength int
	CreatedLength  int
	UpdatedLength  int
}

func NewLengths(lf file.ListFile) *Lengths {
	sl := MaxStringLengths(lf)
	return &Lengths{
		FileNameLength: sl.FileNameLength,
		CreatedLength:  sl.CreatedLength,
		UpdatedLength:  sl.UpdatedLength,
	}
}

func MakeTable(lf file.ListFile) string {
	l := NewLengths(lf)
	upHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftUp, te.RepeatLine(l.FileNameLength), te.MiddleUp, te.RepeatLine(l.CreatedLength), te.MiddleUp, te.RepeatLine(l.UpdatedLength), te.RightUp)
	midHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.V, Fitting(te.tName, l.FileNameLength), te.V, Fitting(te.tCreated, l.CreatedLength), te.V, Fitting(te.tUpdated, l.UpdatedLength), te.V)
	downHead := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftMiddle, te.RepeatLine(l.FileNameLength), te.CenterMiddle, te.RepeatLine(l.CreatedLength), te.CenterMiddle, te.RepeatLine(l.UpdatedLength), te.RightMiddle)
	var table string
	for _, v := range lf {
		table = fmt.Sprintf("%s%c%s%c%s%c%s%c\n", table, te.V, Fitting(v.FileName, l.FileNameLength), te.V, Fitting(v.Created, l.CreatedLength), te.V, Fitting(v.Updated, l.UpdatedLength), te.V)
	}
	footer := fmt.Sprintf("%c%s%c%s%c%s%c", te.LeftBottom, te.RepeatLine(l.FileNameLength), te.MiddleBottom, te.RepeatLine(l.CreatedLength), te.MiddleBottom, te.RepeatLine(l.UpdatedLength), te.RightBottom)
	result := fmt.Sprintf("%s\n%s\n%s\n%s%s", upHead, midHead, downHead, table, footer)
	return result
}

func MaxStringLengths(lf file.ListFile) Lengths {
	var maxFileNameLength, maxCreatedLength, maxUpdatedLength int
	air := 2
	for _, v := range lf {
		if len(v.FileName) > maxFileNameLength {
			maxFileNameLength = len(v.FileName)
		}
		if len(v.Created) > maxCreatedLength {
			maxCreatedLength = len(v.Created)
		}
		if len(v.Updated) > maxUpdatedLength {
			maxUpdatedLength = len(v.Updated)
		}
	}
	if maxFileNameLength < len(te.tName) {
		maxFileNameLength = len(te.tName)
	}
	if maxCreatedLength < len(te.tCreated) {
		maxCreatedLength = len(te.tCreated)
	}
	if maxUpdatedLength < len(te.tUpdated) {
		maxUpdatedLength = len(te.tUpdated)
	}
	return Lengths{maxFileNameLength + air, maxCreatedLength + air, maxUpdatedLength + air}
}

func Fitting(s string, n int) string {
	for len(s) < n {
		s = fmt.Sprintf("%s%c", s, te.WhiteSpace)
		if len(s) == n {
			break
		}
		s = fmt.Sprintf("%c%s", te.WhiteSpace, s)
		if len(s) == n {
			break
		}
	}
	return s
}

type CP struct {
	WhiteSpace   rune
	LeftUp       rune
	MiddleUp     rune
	RightUp      rune
	V            rune
	H            rune
	LeftMiddle   rune
	CenterMiddle rune
	RightMiddle  rune
	LeftBottom   rune
	MiddleBottom rune
	RightBottom  rune
	tName        string
	tCreated     string
	tUpdated     string
}

var te = NewCP()

func (te CP) RepeatLine(n int) string {
	return strings.Repeat(string(te.H), n)
}

func NewCP() *CP {
	return &CP{
		WhiteSpace:   '\u0020',
		LeftUp:       '\u2554',
		MiddleUp:     '\u2566',
		RightUp:      '\u2557',
		V:            '\u2551',
		H:            '\u2550',
		LeftMiddle:   '\u2560',
		CenterMiddle: '\u256c',
		RightMiddle:  '\u2563',
		LeftBottom:   '\u255a',
		MiddleBottom: '\u2569',
		RightBottom:  '\u255d',
		tName:        "File name",
		tCreated:     "Created",
		tUpdated:     "Last updated",
	}
}

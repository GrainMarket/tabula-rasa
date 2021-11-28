package tables

import "fmt"

// ternary is a shim to allow ternary operations in Go
func ternary(check bool, valid interface{}, invalid interface{}) interface{} {
	if check {
		return valid
	}
	return invalid
}

func max(val ...int) (max int) {
	max = val[0]
	for i := 1; i < len(val); i++ {
		if val[i] > max {
			max = val[i]
		}
	}
	return
}

func (tbl *Table) FillWidths() {
	for _, row := range tbl.rows {
		for col, cell := range row {
			tbl.columnWidths[tbl.columns[col]] = max(tbl.columnWidths[tbl.columns[col]], len(cell))
		}
	}
}

func (tbl *Table) CalcWidth(column string, pad bool, verbose bool) (calcWidth int, debug debugCol) {
	i := getSliceIndexString(column, tbl.columns)
	debug.ColName = tbl.columns[i]
	debug.Chars = tbl.CharWidth(debug.ColName)
	debug.PaddingBefore = tbl.Padding(true, i)
	debug.PaddingAfter = tbl.Padding(false, i)

	calcWidth = tbl.CharWidth(debug.ColName)
	if pad {
		calcWidth = tbl.columnWidths[debug.ColName] + tbl.Padding(true, i) + tbl.Padding(false, i)
	}
	if verbose {
		center, _ := tbl.GetBorder(Center)
		fmt.Printf("Column %s(%d) > width: %d + %d + %d =  %d (%t)\n", debug.ColName, i, debug.Chars, debug.PaddingBefore, debug.PaddingAfter, calcWidth, center)
	}

	return
}

func (tbl *Table) Padding(before bool, colIndex int) int {
	if colIndex == 0 { // First column
		if (before && !tbl.borders.showLeft) || (!before && !tbl.borders.showCenter) { // spacing before content with no left border or spacing after content with no center border
			return 0
		}
		return DefaultPadding // spacing after content, or before content with left border

	} else if colIndex == len(tbl.columns)-1 { // Last column
		if before || tbl.borders.showRight { // spacing before content, or after content with right border
			return DefaultPadding
		}
		return 0 // spacing after content with no right border

	} else { // Middle columns
		if tbl.borders.showCenter { // spacing before content, or after content with center border
			return DefaultPadding
		}
		return DefaultPadding // spacing after content with no center border

	}
}

func (tbl *Table) CharWidth(column string) int {
	return tbl.columnWidths[column]
}

func (tbl *Table) ColumnCount() int {
	return len(tbl.columns)
}

func (tbl *Table) GetBorder(border int) (bool, bool) {
	switch border {
	case Top:
		return tbl.borders.showTop, tbl.borders.boldTop
	case Header:
		return tbl.borders.showHeader, tbl.borders.boldHeader
	case Horizontal:
		return tbl.borders.showHorizontal, tbl.borders.boldHorizontal
	case Bottom:
		return tbl.borders.showBottom, tbl.borders.boldBottom
	case Left:
		return tbl.borders.showLeft, tbl.borders.boldLeft
	case Center:
		return tbl.borders.showCenter, tbl.borders.boldCenter
	case Right:
		return tbl.borders.showRight, tbl.borders.boldRight
	}
	return false, false
}

// includesInt checks if the integer is in the array, returning the index if it is, or -1 if it isn't
func getSliceIndexInt(needle int, haystack []int) (index int) {
	var val int
	for index, val = range haystack {
		if val == needle {
			return
		}
	}
	index = -1
	return
}

// includesString checks if the string is in the array, returning the index if it is, or -1 if it isn't
func getSliceIndexString(needle string, haystack []string) (index int) {
	var val string
	for index, val = range haystack {
		if val == needle {
			return
		}
	}
	index = -1
	return
}

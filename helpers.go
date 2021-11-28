package tables

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

func (tbl *Table) calcWidth(column string, pad bool) int {
	if pad {
		i := 0
		for col := range tbl.columnWidths {
			if col == column {
				break
			}
			i++
		}
		return tbl.columnWidths[column] + tbl.padding(true, i) + tbl.padding(false, i)
		// if tbl.borders.showCenter {
		// 	return tbl.columnWidths[column] + (DefaultPadding * 2)
		// }
		// return tbl.columnWidths[column] + DefaultPadding
	}
	return tbl.columnWidths[column]
}

func (tbl *Table) padding(before bool, colIndex int) int {
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

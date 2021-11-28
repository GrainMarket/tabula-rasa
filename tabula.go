package tables

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type Table struct {
	borders         borders
	columns         []string
	columnAlignment map[string]int
	headerAlignment map[string]int
	rows            [][]string
	columnWidths    map[string]int
	// writer
}

type borders struct {
	showTop        bool
	showBottom     bool
	showHeader     bool
	showHorizontal bool
	showLeft       bool
	showCenter     bool
	showRight      bool

	boldTop        bool
	boldBottom     bool
	boldHeader     bool
	boldHorizontal bool
	boldLeft       bool
	boldCenter     bool
	boldRight      bool
}

type WidthFunc func(string) int

const (
	Left = 1 << iota
	Center
	Right
	Top
	Header
	Bottom
	Horizontal
	Left2
	Center2
	Right2
	Top2
	Header2
	Bottom2
	Horizontal2
)

var (
	// DefaultPadding specifies the number of spaces around content in columns.
	DefaultPadding = 1

	// DefaultWriter specifies the output io.Writer for the Table.Print method.
	DefaultWriter io.Writer = os.Stdout

	// DefaultWidthFunc specifies the default WidthFunc for calculating column widths
	DefaultWidthFunc WidthFunc = utf8.RuneCountInString
)

func GetAlignment(align int) (label string) {
	switch align {
	case Left:
		label = "Left"
	case Center:
		label = "Center"
	case Right:
		label = "Right"
	}
	return
}

func NewTable(colNames ...string) (tbl *Table) {
	columnAlignment := make(map[string]int)
	headerAlignment := make(map[string]int)
	columnWidths := make(map[string]int)
	for _, colName := range colNames {
		columnAlignment[colName] = Left
		headerAlignment[colName] = Left
		columnWidths[colName] = len(colName)
	}
	tbl = &Table{
		borders: borders{
			showTop:        false,
			showBottom:     false,
			showHeader:     false,
			showHorizontal: false,
			showLeft:       false,
			showCenter:     false,
			showRight:      false,
			boldTop:        false,
			boldBottom:     false,
			boldHeader:     false,
			boldHorizontal: false,
			boldLeft:       false,
			boldCenter:     false,
			boldRight:      false,
		},
		columns:         colNames,
		columnAlignment: columnAlignment,
		headerAlignment: headerAlignment,
		columnWidths:    columnWidths,
	}

	return
}

func (tbl *Table) SetBorder(border int, display bool, style bool) {
	switch border {
	case Top:
		tbl.borders.showTop = display
		tbl.borders.boldTop = style
	case Header:
		tbl.borders.showHeader = display
		tbl.borders.boldHeader = style
	case Horizontal:
		tbl.borders.showHorizontal = display
		tbl.borders.boldHorizontal = style
	case Bottom:
		tbl.borders.showBottom = display
		tbl.borders.boldBottom = style
	case Left:
		tbl.borders.showLeft = display
		tbl.borders.boldLeft = style
	case Center:
		tbl.borders.showCenter = display
		tbl.borders.boldCenter = style
	case Right:
		tbl.borders.showRight = display
		tbl.borders.boldRight = style
	}
}

func (tbl *Table) AddRow(row ...string) (err error) {
	if len(row) != len(tbl.columns) {
		err = fmt.Errorf("Row length (%d) does not match table columns (%d)", len(row), len(tbl.columns))
		return
	}
	tbl.rows = append(tbl.rows, row)
	return
}

func (tbl *Table) Align(colName string, alignment int, includeHeader bool) {
	tbl.columnAlignment[colName] = alignment
	if includeHeader {
		tbl.headerAlignment[colName] = alignment
	}
}

func (tbl *Table) Print() {
	for _, row := range tbl.rows {
		for col, cell := range row {
			tbl.columnWidths[tbl.columns[col]] = max(tbl.columnWidths[tbl.columns[col]], len(cell))
		}
	}

	tbl.printTopBorder()
	tbl.printHeaders()
	tbl.printHeaderBorder()
	tbl.printRows()
	tbl.printBottomBorder()

	// columnBaseTemplate := "│ %%%dv │ %%-%ds │ %%%ds │ %%-%ds │ %%-%ds │ %%-%ds │ %%-%ds │ %%-%ds │\n"
	// line1BaseTemplate := "┏%s┳%s┳%s┳%s┳%s┳%s┳%s┳%s┓\n"
	// line2BaseTemplate := "┡%s╇%s╇%s╇%s╇%s╇%s╇%s╇%s┩\n"
	// line3BaseTemplate := "└%s┴%s┴%s┴%s┴%s┴%s┴%s┴%s┘\n"
}

func (tbl *Table) printTopBorder() {
	if tbl.borders.showTop {
		if tbl.borders.showLeft {
			if tbl.borders.boldTop && tbl.borders.boldLeft {
				fmt.Print("┏")
			} else if tbl.borders.boldTop && !tbl.borders.boldLeft {
				fmt.Print("┍")
			} else if !tbl.borders.boldTop && tbl.borders.boldLeft {
				fmt.Print("┎")
			} else {
				fmt.Print("┌")
			}
		}
		for i := 0; i < len(tbl.columns); i++ {
			if tbl.borders.boldTop {
				fmt.Print(strings.Repeat("━", tbl.calcWidth(tbl.columns[i], true)))
			} else {
				fmt.Print(strings.Repeat("─", tbl.calcWidth(tbl.columns[i], true)))
			}
			if tbl.borders.showCenter && i < len(tbl.columns)-1 {
				if tbl.borders.boldTop && tbl.borders.boldCenter {
					fmt.Print("┳")
				} else if tbl.borders.boldTop && !tbl.borders.boldCenter {
					fmt.Print("┯")
				} else if !tbl.borders.boldTop && tbl.borders.boldCenter {
					fmt.Print("┰")
				} else {
					fmt.Print("┬")
				}
			}
		}
		if tbl.borders.showRight {
			if tbl.borders.boldTop && tbl.borders.boldRight {
				fmt.Print("┓")
			} else if tbl.borders.boldTop && !tbl.borders.boldRight {
				fmt.Print("┑")
			} else if !tbl.borders.boldTop && tbl.borders.boldRight {
				fmt.Print("┒")
			} else {
				fmt.Print("┐")
			}
		}
		fmt.Print("\n")
	}

}

func (tbl *Table) printHeaderBorder() {
	if tbl.borders.showHeader {
		if tbl.borders.showLeft && tbl.borders.showTop {
			if tbl.borders.boldHeader && tbl.borders.boldLeft {
				fmt.Print("┣")
			} else if tbl.borders.boldHeader && !tbl.borders.boldLeft {
				fmt.Print("┝")
			} else if !tbl.borders.boldHeader && tbl.borders.boldLeft {
				fmt.Print("┠")
			} else {
				fmt.Print("├")
			}
		} else if tbl.borders.showLeft && !tbl.borders.showTop {
			if tbl.borders.boldHeader && tbl.borders.boldLeft {
				fmt.Print("┏")
			} else if tbl.borders.boldHeader && !tbl.borders.boldLeft {
				fmt.Print("┍")
			} else if !tbl.borders.boldHeader && tbl.borders.boldLeft {
				fmt.Print("┎")
			} else {
				fmt.Print("┌")
			}
		}
		for i := 0; i < len(tbl.columns); i++ {
			if tbl.borders.boldHeader {
				fmt.Print(strings.Repeat("━", tbl.calcWidth(tbl.columns[i], true)))
			} else {
				fmt.Print(strings.Repeat("─", tbl.calcWidth(tbl.columns[i], true)))
			}
			if tbl.borders.showCenter && i < len(tbl.columns)-1 {
				if tbl.borders.showCenter { //&& tbl.borders.showTop {
					if tbl.borders.boldHeader && tbl.borders.boldCenter {
						fmt.Print("╋")
					} else if tbl.borders.boldHeader && !tbl.borders.boldCenter {
						fmt.Print("╇")
					} else if !tbl.borders.boldHeader && tbl.borders.boldCenter {
						fmt.Print("╂")
					} else {
						fmt.Print("┼")
					}
					// } else if tbl.borders.showCenter && !tbl.borders.showTop {
					// 	if tbl.borders.boldHeader && tbl.borders.boldCenter {
					// 		fmt.Print("┳")
					// 	} else if tbl.borders.boldHeader && !tbl.borders.boldCenter {
					// 		fmt.Print("┯")
					// 	} else if !tbl.borders.boldHeader && tbl.borders.boldCenter {
					// 		fmt.Print("┰")
					// 	} else {
					// 		fmt.Print("┬")
					// 	}
					// } else if !tbl.borders.showCenter && tbl.borders.showTop {
				} else {
					if tbl.borders.boldHeader && tbl.borders.boldCenter {
						fmt.Print("━")
					} else if tbl.borders.boldHeader && !tbl.borders.boldCenter {
						fmt.Print("%")
					} else if !tbl.borders.boldHeader && tbl.borders.boldCenter {
						fmt.Print("@")
					} else {
						fmt.Print("─")
					}
				}
			}
		}
		if tbl.borders.showRight && tbl.borders.showTop {
			if tbl.borders.boldHeader && tbl.borders.boldRight {
				fmt.Print("┫")
			} else if tbl.borders.boldHeader && !tbl.borders.boldRight {
				fmt.Print("┩")
			} else if !tbl.borders.boldHeader && tbl.borders.boldRight {
				fmt.Print("┨")
			} else {
				fmt.Print("┤")
			}
		} else if tbl.borders.showRight && !tbl.borders.showTop {
			if tbl.borders.boldHeader && tbl.borders.boldRight {
				fmt.Print("┓")
			} else if tbl.borders.boldHeader && !tbl.borders.boldRight {
				fmt.Print("┑")
			} else if !tbl.borders.boldHeader && tbl.borders.boldRight {
				fmt.Print("┒")
			} else {
				fmt.Print("┐")
			}
		}

		fmt.Print("\n")
	}
}

func (tbl *Table) printBottomBorder() {
	if tbl.borders.showBottom {
		if tbl.borders.showLeft {
			if tbl.borders.boldBottom && tbl.borders.boldLeft {
				fmt.Print("┗")
			} else if tbl.borders.boldBottom && !tbl.borders.boldLeft {
				fmt.Print("┕")
			} else if !tbl.borders.boldBottom && tbl.borders.boldLeft {
				fmt.Print("┖")
			} else {
				fmt.Print("└")
			}
		}
		for i := 0; i < len(tbl.columns); i++ {

			if tbl.borders.boldBottom {
				fmt.Print(strings.Repeat("━", tbl.calcWidth(tbl.columns[i], true)))
			} else {
				fmt.Print(strings.Repeat("─", tbl.calcWidth(tbl.columns[i], true)))
			}
			if tbl.borders.showCenter && i < len(tbl.columns)-1 {
				if tbl.borders.boldBottom && tbl.borders.boldCenter {
					fmt.Print("┻")
				} else if tbl.borders.boldBottom && !tbl.borders.boldCenter {
					fmt.Print("┷")
				} else if !tbl.borders.boldBottom && tbl.borders.boldCenter {
					fmt.Print("┸")
				} else {
					fmt.Print("┴")
				}
			}
		}
		if tbl.borders.showRight {
			if tbl.borders.boldBottom && tbl.borders.boldRight {
				fmt.Print("┛")
			} else if tbl.borders.boldBottom && !tbl.borders.boldRight {
				fmt.Print("┙")
			} else if !tbl.borders.boldBottom && tbl.borders.boldRight {
				fmt.Print("┚")
			} else {
				fmt.Print("┘")
			}
			// } else {
			// 	if tbl.borders.boldTop {
			// 		fmt.Print("━")
			// 	} else {
			// 		fmt.Print("─")
			// 	}
		}
		fmt.Print("\n")
	}

}

func (tbl *Table) printHeaders() {
	if tbl.borders.showLeft {
		if tbl.borders.boldLeft {
			fmt.Print("┃")
		} else {
			fmt.Print("│")
		}
	}
	for i := 0; i < len(tbl.columns); i++ {
		fmt.Printf(
			fmt.Sprintf(
				"%s%%%s%ds%s",
				strings.Repeat(" ", tbl.padding(true, i)),
				ternary(tbl.headerAlignment[tbl.columns[i]] == 1, "-", "").(string),
				tbl.columnWidths[tbl.columns[i]],
				strings.Repeat(" ", tbl.padding(false, i)),
			),
			tbl.columns[i],
		)
		if i < len(tbl.columns)-1 {
			if tbl.borders.showCenter {
				if tbl.borders.boldCenter {
					fmt.Print("┃")
				} else {
					fmt.Print("│")
				}
			}
		}
	}
	if tbl.borders.showRight {
		if tbl.borders.boldRight {
			fmt.Print("┃")
		} else {
			fmt.Print("│")
		}
	}
	fmt.Print("\n")
}

func (tbl *Table) printRows() {
	for i := 0; i < len(tbl.rows); i++ {
		tbl.printCells(i)
		if i < len(tbl.rows)-1 {
			tbl.printHorizontal()
		}
	}
}

func (tbl *Table) printCells(rowNum int) {
	if tbl.borders.showLeft {
		if tbl.borders.boldLeft {
			fmt.Print("┃")
		} else {
			fmt.Print("│")
		}
	}
	for i := 0; i < len(tbl.columns); i++ {
		fmt.Printf(
			fmt.Sprintf(
				"%s%%%s%ds%s",
				strings.Repeat(" ", tbl.padding(true, i)),
				ternary(tbl.columnAlignment[tbl.columns[i]] == 1, "-", "").(string),
				tbl.columnWidths[tbl.columns[i]],
				strings.Repeat(" ", tbl.padding(false, i)),
			),
			tbl.rows[rowNum][i],
		)
		if i < len(tbl.columns)-1 {
			if tbl.borders.showCenter {
				if tbl.borders.boldCenter {
					fmt.Print("┃")
				} else {
					fmt.Print("│")
				}
			}
		}
	}

	if tbl.borders.showRight {
		if tbl.borders.boldRight {
			fmt.Print("┃")
		} else {
			fmt.Print("│")
		}
	}
	fmt.Print("\n")
}

func (tbl *Table) printHorizontal() {
	if tbl.borders.showHorizontal {
		if tbl.borders.showLeft {
			if tbl.borders.boldHorizontal && tbl.borders.boldLeft {
				fmt.Print("┣")
			} else if tbl.borders.boldHorizontal && !tbl.borders.boldLeft {
				fmt.Print("┝")
			} else if !tbl.borders.boldHorizontal && tbl.borders.boldLeft {
				fmt.Print("┠")
			} else {
				fmt.Print("├")
			}
		}

		for i := 0; i < len(tbl.columns); i++ {
			if tbl.borders.boldHorizontal {
				fmt.Print(strings.Repeat("━", tbl.calcWidth(tbl.columns[i], true)))
			} else {
				fmt.Print(strings.Repeat("─", tbl.calcWidth(tbl.columns[i], true)))
			}
			if tbl.borders.showCenter && i < len(tbl.columns)-1 {
				if tbl.borders.boldHorizontal && tbl.borders.boldCenter {
					fmt.Print("╋")
				} else if tbl.borders.boldHorizontal && !tbl.borders.boldCenter {
					fmt.Print("┿")
				} else if !tbl.borders.boldHorizontal && tbl.borders.boldCenter {
					fmt.Print("╂")
				} else {
					fmt.Print("┼")
				}
			}
		}
		if tbl.borders.showRight {
			if tbl.borders.boldHorizontal && tbl.borders.boldRight {
				fmt.Print("┫")
			} else if tbl.borders.boldHorizontal && !tbl.borders.boldRight {
				fmt.Print("┥")
			} else if !tbl.borders.boldHorizontal && tbl.borders.boldRight {
				fmt.Print("┨")
			} else {
				fmt.Print("┤")
			}
		}
		fmt.Print("\n")
	}

}

package main

import (
	"fmt"
	"strconv"
	"strings"

	tabula "github.com/GrainMarket/tabula-rasa"
)

func demo() {
	columns := []string{"Col1", "Col2", "Col3", "Col4"}
	tbl := tabula.NewTable(columns...)
	err := tbl.AddRow("Something longer than the column header", "short", "3.14", "")
	if err != nil {
		fmt.Print("Expected no error, got", err)
	}
	err = tbl.AddRow("x", "y")
	if err == nil {
		fmt.Print("Expected error, got nothing")
	}
	for i := 0; i < 3; i++ {
		err := tbl.AddRow(strconv.Itoa(i+1), strconv.Itoa((i+1)*2), strconv.Itoa((i+1)*3), strconv.Itoa((i+1)*4))
		if err != nil {
			fmt.Print("Expected no error, got", err)
		}
	}

	tbl.SetBorder(tabula.Center, true, false)
	tbl.SetBorder(tabula.Header, true, true)
	tbl.SetBorder(tabula.Top, true, false)

	tbl.fillWidths()
	for i := 0; i < tbl.ColumnCount(); i++ {
		padding := tbl.padding(true, i)
		width := tbl.calcWidth(columns[i], true)
		fmt.Printf("Column %d > width = %d | pre-content padding = %d | ", i, width, padding)
		padding = tbl.padding(false, i)
		fmt.Printf("post-content padding = %d (%t)\n", padding, tbl.borders.showCenter)
	}
	for i := 0; i < tbl.ColumnCount(); i++ {
		if i != tbl.ColumnCount()-1 {
			fmt.Printf("%s|", strings.Repeat("_", tbl.calcWidth(columns[i], true)))
		} else {
			fmt.Printf("%s\n", strings.Repeat("-", tbl.calcWidth(columns[i], true)))
		}
	}

	tbl.Print()
}

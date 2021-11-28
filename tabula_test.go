package tables

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestOneTable(t *testing.T) {
	tbl := NewTable("Col1", "Col2", "Col3", "Col4")
	if len(tbl.columns) != 4 {
		t.Error("Expected 2 columns, got", len(tbl.columns))
	}
	err := tbl.AddRow("Something longer than the column header", "short", "3.14", "")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	err = tbl.AddRow("x", "y")
	if err == nil {
		t.Error("Expected error, got nothing")
	}
	for i := 0; i < 3; i++ {
		err := tbl.AddRow(strconv.Itoa(i+1), strconv.Itoa((i+1)*2), strconv.Itoa((i+1)*3), strconv.Itoa((i+1)*4))
		if err != nil {
			t.Error("Expected no error, got", err)
		}
	}

	tbl.SetBorder(Center, true, false)
	tbl.SetBorder(Header, true, false)
	// tbl.SetBorder(Top, true, false)

	for i := 0; i < len(tbl.columns); i++ {
		padding := tbl.padding(true, i)
		fmt.Printf("Column %d pre-content padding = %d (%t)\n", i, padding, tbl.borders.showCenter)
		padding = tbl.padding(false, i)
		fmt.Printf("Column %d post-content padding = %d (%t)\n", i, padding, tbl.borders.showCenter)
	}

	tbl.Print()
}

func TestNewTable(t *testing.T) {
	tbl := NewTable("Col1", "Col2", "Col 3")
	if len(tbl.columns) != 3 {
		t.Error("Expected 2 columns, got", len(tbl.columns))
	}
	err := tbl.AddRow("a", "b", "c")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	err = tbl.AddRow("x", "y")
	if err == nil {
		t.Error("Expected error, got nothing")
	}
	err = tbl.AddRow("d", "e", "f")
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	err = tbl.AddRow("g", "h", "i")
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	tbl.Align("Col1", Right, true)
	tbl.Align("Col1", Left, false)
	tbl.Align("Col3", Right, true)
	tbl.Align("Col2", Left, true)
	tbl.Align("Col2", Right, false)

	if tbl.columnAlignment["Col1"] != Left {
		t.Error("Column 1 Alignment - Expected Left, got ", GetAlignment(tbl.columnAlignment["Col1"]))
	}
	if tbl.columnAlignment["Col2"] != Right {
		t.Error("Column 2 Alignment - Expected Right, got ", GetAlignment(tbl.columnAlignment["Col1"]))
	}
	if tbl.headerAlignment["Col1"] != Right {
		t.Error("Header 1 Alignment - Expected Right, got ", GetAlignment(tbl.columnAlignment["Col1"]))
	}
	if tbl.headerAlignment["Col2"] != Left {
		t.Error("Header 2 Alignment - Expected Left, got ", GetAlignment(tbl.columnAlignment["Col1"]))
	}

	// 6 3 1
	// 4	2	6	8	4	2	1
	// -------------
	// 0 0 0 0 0 0 0 = 0
	// 1 1 1 1 1 1 1 = 127
	// borders	:= []string{ "Left","Center","Right","Top","Header","Bottom","Horizontal"}

	// fmt.Printf("%d %d %d %d %d %d %d\n", Left, Center, Right, Top, Header, Bottom, Horizontal)

	// test := buildConf(false, true, true, false, false, false, false)
	// bw := test & Left
	// fmt.Printf("%d %s %s (bw = %d)", test, Ternary(bw != 0, "includes", "does not include").(string), GetAlignment(Left), bw)

	// test = Left + Center
	// match = test | Left
	// fmt.Printf("%d vs %d ", match, test)

	// test = Left + Right
	// match = test | Left
	// fmt.Printf("%d vs %d ", match, test)
	line := 0
	var done []int
	var next int

	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			next = buildConf((j&Left != 0), (j&Center != 0), (j&Right != 0), (j&Top != 0), (j&Header != 0), (j&Bottom != 0), (j&Horizontal != 0), (i&Left != 0), (i&Center != 0), (i&Right != 0), (i&Top != 0), (i&Header != 0), (i&Bottom != 0), (i&Horizontal != 0))
			if intInArray(next, done) == -1 {
				done = append(done, next)
				printWithConf(tbl, (j&Left != 0), (j&Center != 0), (j&Right != 0), (j&Top != 0), (j&Header != 0), (j&Bottom != 0), (j&Horizontal != 0), (i&Left != 0), (i&Center != 0), (i&Right != 0), (i&Top != 0), (i&Header != 0), (i&Bottom != 0), (i&Horizontal != 0))
				if line%3 == 2 {
					time.Sleep(1 * time.Second)
				}
				line++
			}
		}
	}

}

func buildConf(sLeft, sCenter, sRight, sTop, sHeader, sBottom, sHorizonal bool, dLeft, dCenter, dRight, dTop, dHeader, dBottom, dHorizonal bool) (displayConf int) {
	if dLeft {
		displayConf = displayConf + Left
		if sLeft {
			displayConf = displayConf + Left2
		}
	}
	if dCenter {
		displayConf = displayConf + Center
		if sCenter {
			displayConf = displayConf + Center2
		}
	}
	if dRight {
		displayConf = displayConf + Right
		if sRight {
			displayConf = displayConf + Right2
		}
	}
	if dTop {
		displayConf = displayConf + Top
		if sTop {
			displayConf = displayConf + Top2
		}
	}
	if dHeader {
		displayConf = displayConf + Header
		if sHeader {
			displayConf = displayConf + Header2
		}
	}
	if dBottom {
		displayConf = displayConf + Bottom
		if sBottom {
			displayConf = displayConf + Bottom2
		}
	}
	if dHorizonal {
		displayConf = displayConf + Horizontal
		if sHorizonal {
			displayConf = displayConf + Horizontal2
		}
	}
	return
}

func printWithConf(tbl *Table, sLeft, sCenter, sRight, sTop, sHeader, sBottom, sHorizontal bool, dLeft, dCenter, dRight, dTop, dHeader, dBottom, dHorizontal bool) {
	tbl.SetBorder(Header, false, false)
	tbl.SetBorder(Bottom, false, false)
	tbl.SetBorder(Horizontal, false, false)
	label := ""

	if dLeft {
		label += "Left"
		tbl.SetBorder(Left, true, false)
		if sLeft {
			label += "+Bold"
			tbl.SetBorder(Left, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Left, false, false)
	}

	if dCenter {
		label += "Center"
		tbl.SetBorder(Center, true, false)
		if sCenter {
			label += "+Bold"
			tbl.SetBorder(Center, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Center, false, false)
	}

	if dRight {
		label += "Right"
		tbl.SetBorder(Right, true, false)
		if sRight {
			label += "+Bold"
			tbl.SetBorder(Right, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Right, false, false)
	}

	if dTop {
		label += "Top"
		tbl.SetBorder(Top, true, false)
		if sTop {
			label += "+Bold"
			tbl.SetBorder(Top, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Top, false, false)
	}

	if dHeader {
		label += "Header"
		tbl.SetBorder(Header, true, false)
		if sHeader {
			label += "+Bold"
			tbl.SetBorder(Header, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Header, false, false)
	}

	if dBottom {
		label += "Bottom"
		tbl.SetBorder(Bottom, true, false)
		if sBottom {
			label += "+Bold"
			tbl.SetBorder(Bottom, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Bottom, false, false)
	}

	if dHorizontal {
		label += "Horizontal"
		tbl.SetBorder(Horizontal, true, false)
		if sHorizontal {
			label += "+Bold"
			tbl.SetBorder(Horizontal, true, true)
		}
		label += ", "
	} else {
		tbl.SetBorder(Horizontal, false, false)
	}

	fmt.Println(label)
	tbl.Print()
	fmt.Println()
}

// Header, Vertical (HeaderBorder Missing)
// Header, Vertical+Bold (HeaderBorder Missing)
// Header+Bold, Vertical (HeaderBorder Missing)
// Header+Bold, Vertical+Bold (HeaderBorder Missing)
// Top, Header+Bold, Vertical (HeaderBorder Bold Vertical on middle & end)
// Top+Bold, Header+Bold, Vertical (HeaderBorder Bold Vertical on middle & end)

// intInArray checks if the integer is in the array
func intInArray(needle int, haystack []int) (index int) {
	var val int
	for index, val = range haystack {
		if val == needle {
			return
		}
	}
	index = -1
	return
}

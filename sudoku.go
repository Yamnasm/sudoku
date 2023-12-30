package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	TEST_PUZZLE   = "920300154050001000300042700040000075009000200670000090006920008000600040197005023"
	COLOUR_RESET  = "\033[0m"
	COLOUR_RED    = "\033[31m"
	COLOUR_GREEN  = "\033[32m"
	COLOUR_YELLOW = "\033[33m"
	COLOUR_BLUE   = "\033[34m"
	COLOUR_PURPLE = "\033[35m"
	COLOUR_CYAN   = "\033[36m"
	COLOUR_GREY   = "\033[37m"
	COLOUR_WHITE  = "\033[97m"
	COLOUR_BLACK  = "\033[30m"
)

type Cell struct {
	value       string
	potentials  []string
	coordinates Position
}

type Position struct {
	row    int
	column int
	box    int
}

func main() {
	colour_init()

	puz := parse_puzzle_string(TEST_PUZZLE)

	fmt.Println(get_entire_box(puz, 0))

	// small_display(TEST_PUZZLE)
	// large_display(TEST_PUZZLE)

}

func colour_init() {
	// Opt-in for ansi color support for current process.
	// https://learn.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences#output-sequences
	var outMode uint32
	out := windows.Handle(os.Stdout.Fd())
	if err := windows.GetConsoleMode(out, &outMode); err != nil {
		return
	}
	outMode |= windows.ENABLE_PROCESSED_OUTPUT | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	_ = windows.SetConsoleMode(out, outMode)
}

func small_display(puz string) {
	for i := 1; i < 10; i++ {
		var row string
		row = ""
		for x := 1; x < 10; x++ {
			index := (i * x) - 1

			if x-1 > 0 && (x-1)%3 == 0 {
				row = row + "|"
			}

			if string(puz[index]) == "0" {
				row = row + COLOUR_WHITE + " " + string(puz[index]) + " " + COLOUR_RESET
			} else {
				row = row + COLOUR_GREEN + " " + string(puz[index]) + " " + COLOUR_RESET
			}
		}
		if i-1 > 0 && (i-1)%3 == 0 {
			fmt.Println("─────────┼─────────┼─────────")
		}
		fmt.Println(row)
	}
}

func parse_puzzle_string(puz string) []Cell {
	var cell_array []Cell
	for i, c := range strings.Split(TEST_PUZZLE, "") {
		row := i / 9
		col := i % 9
		box := ((row / 3) * 3) + ((col / 3) % 3)

		p := Cell{value: c, coordinates: Position{row: row, column: col, box: box}}

		cell_array = append(cell_array, p)
	}
	return cell_array
}

// func fill_initial_potentials(puz []Cell) []Cell {
// 	for _, c := range puz {
// 		if c.value != "0" {
// 			continue
// 		}

// 	}
// }

func get_entire_row(puz []Cell, target int) []string {
	var output []string
	for _, c := range puz {
		if c.coordinates.row == target {
			output = append(output, c.value)
		}
	}
	return output
}

func get_entire_col(puz []Cell, target int) []string {
	var output []string
	for _, c := range puz {
		if c.coordinates.column == target {
			output = append(output, c.value)
		}
	}
	return output
}

func get_entire_box(puz []Cell, target int) []string {
	var output []string
	for _, c := range puz {
		if c.coordinates.box == target {
			output = append(output, c.value)
		}
	}
	return output
}

func get_pos(puz []Cell, target []int) {}

func large_display(puz string) {}

/*

      | ∙ ∙ ∙ | 1 2 3 ║
  5   | ∙ 5 ∙ | 4 5 6 ║
      | ∙ ∙ ∙ | 7 8 9 ║
-----+
*/

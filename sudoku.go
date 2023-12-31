package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	TEST_PUZZLE = "457869003009200058000000740000730000620908075000046000061000000730005100500193867"
	LOOP_LIMIT  = 50

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
	small_display(puz)
	solve_loop(puz)
}

func solve_loop(puz []Cell) {

	for cycles := 0; cycles < LOOP_LIMIT; cycles++ {
		if check_complete(puz) {
			fmt.Println("")
			fmt.Println("Sudoku puzzle is complete after", cycles, "cycles!")
			return
		}
		// saving for difference check later
		pre_check_puz := puz

		// running though each solver function
		/* TODO:
		evaluate if assess_potentials should be run after each solve
		function or at the start of the loop.
		*/
		puz = assess_potentials(puz)
		puz = naked_singles(puz)
		// further check strategies
		//puz = hidden_singles(puz)
		//puz = naked_pairs(puz)
		//puz = x_wing(puz)

		fmt.Println("")
		small_display(puz)

		if is_puzzle_same(pre_check_puz, puz) {
			fmt.Println("")
			fmt.Println("Sudoku puzzle stalled after", cycles, "cycles.")
			return
		}
	}
	fmt.Println("")
	fmt.Println("ERROR: Solver hit limit after", LOOP_LIMIT, "cycles.")
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

func small_display(puz []Cell) {
	for i := 0; i < 9; i++ {
		var print_row string
		for p, n := range strings.Join(get_entire_row(puz, i), " ") {
			if string(n) == "0" {
				print_row = print_row + COLOUR_WHITE + string(n) + COLOUR_RESET
			} else {
				print_row = print_row + COLOUR_GREEN + string(n) + COLOUR_RESET
			}
			if (p+1)%6 == 0 {
				print_row = print_row + "| "
			}
		}
		if i != 0 && i%3 == 0 {
			fmt.Println(" ──────┼───────┼──────")
		}
		fmt.Println(" " + print_row + " ")
	}
}

func naked_singles(puz []Cell) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.value == "0" {
			if len(c.potentials) == 1 {
				c.value = c.potentials[0]
			}
		}
		output = append(output, c)
	}
	return output
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

func assess_potentials(puz []Cell) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.value == "0" {
			var potentials []string
			for i := 1; i < 10; i++ {
				potentials = append(potentials, fmt.Sprint(i))
			}
			check_row := get_entire_row(puz, c.coordinates.row)
			potentials = diff(potentials, check_row)

			check_col := get_entire_col(puz, c.coordinates.column)
			potentials = diff(potentials, check_col)

			check_box := get_entire_box(puz, c.coordinates.box)
			potentials = diff(potentials, check_box)

			c.potentials = potentials
		}
		output = append(output, c)
	}
	return output
}

func check_complete(puz []Cell) bool {
	for _, c := range puz {
		if c.value == "0" {
			return false
		}
	}
	return true
}

func is_puzzle_same(previous_puz []Cell, current_puz []Cell) bool {
	return reflect.DeepEqual(previous_puz, current_puz)
}

func diff(a []string, b []string) []string {
	// Turn b into a map
	m := make(map[string]bool, len(b))
	for _, s := range b {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []string
	for _, s := range a {
		if _, ok := m[s]; !ok {
			diff = append(diff, s)
			continue
		}
		m[s] = true
	}
	// Sort the resulting slice
	sort.Strings(diff)
	return diff
}

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

func large_display(puz string) {

	/*
		something like this:

			| ∙ ∙ ∙ | 1 2 3 ║
		5   | ∙ 5 ∙ | 4 5 6 ║
			| ∙ ∙ ∙ | 7 8 9 ║
		-----+
	*/
}

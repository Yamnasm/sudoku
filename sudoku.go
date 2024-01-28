package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	TEST_PUZZLE_EASY   = "457869003009200058000000740000730000620908075000046000061000000730005100500193867"
	TEST_PUZZLE_MEDIUM = "940030010301840005000002070700000000020465090000000002060300000800054607030010048"
	TEST_PUZZLE_HARD   = "000025670609000010002004000080002700030060090004500020000200800020000301096480000"
	TEST_PUZZLE_EVIL   = "003000001004086009000100030030900140800000005027001090070005000300490600200000500"

	LOOP_LIMIT = 100

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

// consts for representing the houses parameter on various functions
const (
	HOUSE_ALL = iota
	HOUSE_ROW
	HOUSE_COL
	HOUSE_BOX
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
	var puz []Cell = parse_puzzle_string(TEST_PUZZLE_EVIL)
	small_display(puz)
	puz = assess_potentials(puz)
	solve_loop(puz)
}

// Loop to apply check/solve functions in a cycle
func solve_loop(puz []Cell) {
	for cycles := 0; cycles < LOOP_LIMIT; cycles++ {
		if check_complete(puz) {
			fmt.Println("")
			fmt.Println("Sudoku puzzle is complete after", cycles, "cycles!")
			small_display(puz)
			return
		}

		// checks if puzzle solver has stalled if all checks are exhausted.
		// TODO: make "Status" struct to cover other check outputs that may
		// come up.
		var change_made bool = false

		// Steps though solution using more complicated checks as each
		// previous check is exhausted. There is probably a better way
		// to do this.
		puz, change_made = open_singles(puz)
		if change_made {
			continue
		}
		puz, change_made = hidden_singles(puz)
		if change_made {
			continue
		}

		// TODO: Add more advanced checks.
		// puz = open_pairs(puz)
		// puz = x_wing(puz)

		if !change_made {
			fmt.Println("")
			large_display(puz)
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

// solves for single potentials within cells
func open_singles(puz []Cell) ([]Cell, bool) {
	change_made := false
	for i, c := range puz {
		if c.value == "0" {
			if len(c.potentials) == 1 {
				fmt.Println("Open Single", c.potentials[0], "found on", c.coordinates.row+1, c.coordinates.column+1)
				c.value = c.potentials[0]
				c.potentials = nil
				puz[i] = c
				puz = eliminate_potentials(c.value, c.coordinates, HOUSE_ALL, puz)
				change_made = true
				return puz, change_made
			}
		}
	}
	return puz, change_made
}

// solves for singular potentials within houses
func hidden_singles(puz []Cell) ([]Cell, bool) {
	change_made := false

	for i, c := range puz {
		if c.value == "0" {
			cell_pos := []int{c.coordinates.row, c.coordinates.column, c.coordinates.box}
			for house := 1; house < 4; house++ {
				entire_house := get_entire_house(puz, cell_pos[house-1], house)
				flat_potentials_list := concatinate_potentials(entire_house)
				for _, p := range c.potentials {
					if count(p, flat_potentials_list) == 1 {
						fmt.Println("Hidden Single", p, "found on", c.coordinates.row+1, c.coordinates.column+1)
						c.value = p
						c.potentials = nil
						puz[i] = c
						puz = eliminate_potentials(c.value, c.coordinates, HOUSE_ALL, puz)
						change_made = true
						return puz, change_made
					}
				}
			}
		}
	}
	return puz, change_made
}

// UNFINISHED identifies twins of potential pairs
func open_pairs(puz []Cell) []Cell {
	var output []Cell
	for _, c := range puz {
		if len(c.potentials) == 2 { //finding the initial pair

		}
		output = append(output, c)
	}
	return output
}

// func check_duplicate_in_arr(target []string, house []Cell) int {
// 	// target is typically an array of potentials. may change this later
// 	for p, c := range house {
// 		if target == c.potentials {
// 			return p
// 		}
// 	}
// }

// counts the target string within a string array
func count(target string, array []string) int {
	var count int = 0
	for _, e := range array {
		if target == e {
			count++
		}
	}
	return count
}

// converts a house of cells into a flat string array of it's potentials
func concatinate_potentials(cells []Cell) []string {
	var joined []string
	var potentials_array [][]string

	for _, p := range cells {
		potentials_array = append(potentials_array, p.potentials)
	}

	for _, s := range potentials_array {
		joined = append(joined, s...)
	}
	return joined
}

// converts initial puzzle string (from random_sudoku.py) to Cell array.
func parse_puzzle_string(puz string) []Cell {
	var cell_array []Cell
	for i, c := range strings.Split(puz, "") {
		row := i / 9
		col := i % 9
		box := ((row / 3) * 3) + ((col / 3) % 3)

		p := Cell{value: c, coordinates: Position{row: row, column: col, box: box}}

		cell_array = append(cell_array, p)
	}
	return cell_array
}

// fills all of the potentials of each cell by the values of houses it belongs to.
// Care must be taken because this func will rewrite previously eliminated potentials
func assess_potentials(puz []Cell) []Cell {
	for i, c := range puz {
		if c.value == "0" {
			var potentials []string
			for i := 1; i < 10; i++ {
				potentials = append(potentials, fmt.Sprint(i))
			}
			cell_pos := []int{c.coordinates.row, c.coordinates.column, c.coordinates.box}
			for house := 1; house < 4; house++ {
				entire_house := get_entire_house(puz, cell_pos[house-1], house)
				potentials = diff(potentials, stringify_cells(entire_house))
			}
			puz[i].potentials = potentials
		}
	}
	return puz
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

// Removes potentials from puzzle based on target string.
// May need a target house argument on pair-checks (and more).
func eliminate_potentials(target string, pos Position, house int, puz []Cell) []Cell {
	for i, c := range puz {

		// mode parameter to remove potentials from specific (or all) houses.
		if house == HOUSE_ALL {
			if c.coordinates.row != pos.row &&
				c.coordinates.column != pos.column &&
				c.coordinates.box != pos.box {
				continue
			}
		}
		if house == HOUSE_ROW {
			if c.coordinates.row != pos.row {
				continue
			}
		}
		if house == HOUSE_COL {
			if c.coordinates.column != pos.column {
				continue
			}
		}
		if house == HOUSE_BOX {
			if c.coordinates.box != pos.box {
				continue
			}
		}

		for _, p := range c.potentials {
			if p == target {
				fmt.Println("Removed", p, "from", c.coordinates.row+1, c.coordinates.column+1)
				c.potentials = remove(c.potentials, target)
			}
		}
		puz[i] = c
	}
	return puz
}

// removes string from string array
func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// converting a house of cell values into a string array
func stringify_cells(cell_arr []Cell) []string {
	var output []string
	for _, c := range cell_arr {
		output = append(output, c.value)
	}
	return output
}

// super quick check if puzzle is not filled in, then long form validation
func check_complete(puz []Cell) bool {

	for _, c := range puz {
		if c.value == "0" {
			return false
		}
	}
	for house := 1; house < 4; house++ {
		// house iterates down the HOUSE_X constants.
		for i := 0; i < 9; i++ {
			entire_house := stringify_cells(get_entire_house(puz, i, house))
			for num := 1; num < 10; num++ {
				if count(fmt.Sprint(num), entire_house) > 1 {
					fmt.Println("Puzzle is invalid!")
					return false
				}
			}
		}
	}
	// if we got to this point, it's all good!
	return true
}

// Returns a string array of values from a selected house.
func get_entire_house(puz []Cell, target int, house int) []Cell {
	var output []Cell
	for _, c := range puz {
		if house == HOUSE_ROW {
			if c.coordinates.row == target {
				output = append(output, c)
			}
		}
		if house == HOUSE_COL {
			if c.coordinates.column == target {
				output = append(output, c)
			}
		}
		if house == HOUSE_BOX {
			if c.coordinates.box == target {
				output = append(output, c)
			}
		}

	}
	return output
}

// prints 9x9 layout of sudoku puzzle showing just values.
func small_display(puz []Cell) {
	for i := 0; i < 9; i++ {
		var print_row string
		entire_row := get_entire_house(puz, i, HOUSE_ROW)
		for p, n := range strings.Join(stringify_cells(entire_row), " ") {
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

// prints 27x27 layout of sudoku puzzle, showing all values and potentials
func large_display(puz []Cell) {

	for i := 0; i < 9; i++ {
		var print_row string
		entire_row := get_entire_house(puz, i, HOUSE_ROW)

		for r := 0; r < 3; r++ {
			for p, n := range entire_row {
				for k := 1; k < 4; k++ {
					check_num := (r * 3) + k
					if n.value != "0" {
						if check_num == 5 {
							print_row = print_row + n.value
							continue
						}
						print_row = print_row + "■"
						continue
					}
					if slices.Contains(n.potentials, fmt.Sprint(check_num)) {
						print_row = print_row + fmt.Sprint(check_num)
						continue
					}
					print_row = print_row + " "
					continue
				}
				if p != 8 {
					print_row = print_row + "|"
				}
			}
			// adding spaces between the values for better readability
			print_row = strings.Join(strings.Split(print_row, ""), " ")
			fmt.Println(" " + print_row)
			print_row = ""

			// WIP:
			// cooler and better lines for beautiful CLI effect
			if r == 2 && i != 8 {
				fmt.Println("───────┼───────┼───────╬───────┼───────┼───────╬───────┼───────┼───────")
			}
		}

		// random unicode for visual candy:
		// ─┼─
		// ■
		// ╬
		// ╫  ╪
	}
}

package main

import (
	"fmt"
	"os"
	"reflect"
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

	LOOP_LIMIT = 50

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
	var puz []Cell = parse_puzzle_string(TEST_PUZZLE_HARD)
	// small_display(puz)
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
		var pre_check_puz []Cell = puz

		// running though each solver function
		/* TODO:
		evaluate if assess_potentials should be run after each solve
		function or at the start of the loop.
		*/
		puz = assess_potentials(puz)
		small_display(puz)
		fmt.Println("")
		large_display(puz)
		// puz = open_singles(puz)
		// puz = assess_potentials(puz)
		// puz = hidden_singles(puz, 1) //mode 1 is rows
		// puz = assess_potentials(puz)
		// puz = hidden_singles(puz, 2) //mode 2 is columns
		// puz = assess_potentials(puz)
		// puz = hidden_singles(puz, 3) //mode 3 is boxes
		//puz = open_pairs(puz)
		//puz = x_wing(puz)

		fmt.Println("")
		// small_display(puz)

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

func open_singles(puz []Cell) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.value == "0" {
			if len(c.potentials) == 1 {
				fmt.Println("Open Single", c.potentials[0])
				c.value = c.potentials[0]
				c.potentials = nil
			}
		}
		output = append(output, c)
	}
	return output
}

// mild attempt at this check function. Doesn't work but can do.
func hidden_singles(puz []Cell, mode int) []Cell {
	var output []Cell

	// checking rows
	if mode == 1 {
		for _, c := range puz {
			if c.value == "0" {
				entire_row := get_entire_row(puz, c.coordinates.row)
				flat_potentials_list := concatinate_potentials(entire_row)
				for _, p := range c.potentials {
					if count(p, flat_potentials_list) == 1 {
						fmt.Println("Hidden Single", p)
						c.value = p
						c.potentials = nil
					}
				}
			}
			output = append(output, c)
		}
	}
	// checking columns
	if mode == 2 {
		for _, c := range puz {
			if c.value == "0" {
				entire_row := get_entire_row(puz, c.coordinates.row)
				flat_potentials_list := concatinate_potentials(entire_row)
				for _, p := range c.potentials {
					if count(p, flat_potentials_list) == 1 {
						fmt.Println("Hidden Single", p)
						c.value = p
						c.potentials = nil
					}
				}
			}
			output = append(output, c)
		}
	}
	// checking boxes
	if mode == 3 {
		for _, c := range puz {
			if c.value == "0" {
				entire_row := get_entire_row(puz, c.coordinates.row)
				flat_potentials_list := concatinate_potentials(entire_row)
				for _, p := range c.potentials {
					if count(p, flat_potentials_list) == 1 {
						fmt.Println("Hidden Single", p)
						c.value = p
						c.potentials = nil
					}
				}
			}
			output = append(output, c)
		}
	}
	return output
}

func count(target string, array []string) int {
	var count int = 0
	for _, e := range array {
		if target == e {
			count++
		}
	}
	return count
}

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

func listify_values(cells []Cell) []string {
	var joined []string

	for _, s := range cells {
		joined = append(joined, s.value)
	}
	return joined
}

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

func assess_potentials(puz []Cell) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.value == "0" {
			var potentials []string
			for i := 1; i < 10; i++ {
				potentials = append(potentials, fmt.Sprint(i))
			}
			entire_row := get_entire_row(puz, c.coordinates.row)
			potentials = diff(potentials, stringify_cells(entire_row))

			entire_col := get_entire_col(puz, c.coordinates.column)
			potentials = diff(potentials, stringify_cells(entire_col))

			entire_box := get_entire_box(puz, c.coordinates.box)
			potentials = diff(potentials, stringify_cells(entire_box))

			c.potentials = potentials
		}
		output = append(output, c)
	}
	return output
}

func stringify_cells(cell_arr []Cell) []string {
	var output []string
	for _, c := range cell_arr {
		output = append(output, c.value)
	}
	return output
}

func check_complete(puz []Cell) bool {
	// super quick check if puzzle is not filled in
	for _, c := range puz {
		if c.value == "0" {
			return false
		}
	}
	// proper confirmation if puzzle is *valid*
	for i := 0; i < 9; i++ {
		entire_row := listify_values(get_entire_row(puz, i))
		entire_col := listify_values(get_entire_col(puz, i))
		entire_box := listify_values(get_entire_box(puz, i))
		for num := 1; num < 10; num++ {
			if count(fmt.Sprint(num), entire_row) > 1 {
				fmt.Println("Puzzle is invalid!")
				return false
			}
			if count(fmt.Sprint(num), entire_col) > 1 {
				fmt.Println("Puzzle is invalid!")
				return false
			}
			if count(fmt.Sprint(num), entire_box) > 1 {
				fmt.Println("Puzzle is invalid!")
				return false
			}
		}
	}
	// if we got to this point, it's all good!
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

func get_entire_row(puz []Cell, target int) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.coordinates.row == target {
			output = append(output, c)
		}
	}
	return output
}

func get_entire_col(puz []Cell, target int) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.coordinates.column == target {
			output = append(output, c)
		}
	}
	return output
}

func get_entire_box(puz []Cell, target int) []Cell {
	var output []Cell
	for _, c := range puz {
		if c.coordinates.box == target {
			output = append(output, c)
		}
	}
	return output
}

func small_display(puz []Cell) {
	for i := 0; i < 9; i++ {
		var print_row string
		entire_row := get_entire_row(puz, i)
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

// working on making the potentials visible on CLI. WIP
func large_display(puz []Cell) {

	for i := 0; i < 9; i++ {
		var print_row string
		entire_row := get_entire_row(puz, i)

		for r := 0; r < 3; r++ {
			for p, n := range entire_row {
				for k := 1; k < 4; k++ {
					check_num := (r * 3) + k
					if slices.Contains(n.potentials, fmt.Sprint(check_num)) {
						print_row = print_row + fmt.Sprint(check_num)
					} else {
						print_row = print_row + " "
					}
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

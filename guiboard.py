from tkinter import *
from SudokuPuzzle import SudokuPuzzle
# from dataclasses import dataclass

BOARD_DIMS = 400
CELL_DIMS = BOARD_DIMS / 9
FONT_SIZE = int(CELL_DIMS / 1.4)

class GUIBoard:
    def __init__(self):
        self.window = Tk()
        self.window.title("Sudoku")
        self.window.resizable(False, False)
        self.window.configure(bg="white")
        self.window.bind("<1>", self._click_canvas)
        self.window.bind("<KeyPress>", self._key_pressed)

        self.board_container = Frame(
            self.window,
            background="white"
        )

        self.canvas = Canvas(
            master = self.board_container,
            height = BOARD_DIMS,
            width  = BOARD_DIMS
        )

        self.cell_rect = []
        self._draw_cells()

        self.box_rect = []
        self._draw_boxes()

        # list of drawn numbers on the canvas.
        self.num_text = []
        self.guess_text = []

        self.selected_cell = None

        export_button = Button(self.window, text="Export", command=self.export)
        export_button.pack()

        self.board_container.pack()
        self.canvas.pack()

    def _draw_boxes(self):
        rect_size = BOARD_DIMS / 3
        for i in range(9):
            x0 = 0 + (rect_size * (i % 3))
            y0 = 0 + (rect_size * (i // 3))
            x1 = rect_size + (rect_size * (i % 3))
            y1 = rect_size + (rect_size * (i // 3))
            self.box_rect.append(
                self.canvas.create_rectangle(
                    x0, y0, x1, y1,
                    width=3
                )
            )
    
    def _draw_cells(self):
        rect_size = BOARD_DIMS / 9
        for i in range(81):
            x0 = 0 + (rect_size * (i % 9))
            y0 = 0 + (rect_size * (i // 9))
            x1 = rect_size + (rect_size * (i % 9))
            y1 = rect_size + (rect_size * (i // 9))
            self.cell_rect.append(
                self.canvas.create_rectangle(
                    x0, y0, x1, y1,
                    fill="white",
                    width=1
                )
            )

    def _update_gui(self, puzzle):

        # removing all visible text from the board.
        for num in self.num_text:
            self.canvas.delete(num)

        for guess in self.guess_text:
            self.canvas.delete(guess)

        self.guess_text = []
        self.num_text = []

        # replacing the text with the "new" puzzle object
        for cell in puzzle:
            if cell.value == "0":
                continue
            row = cell.location.row
            col = cell.location.column
            self.place_number(row, col, cell.value)
    
    def load_puzzle(self, puzzle):
        self.puzzle = puzzle
        self._update_gui(puzzle)

    def place_number(self, row_num, col_num, value):

        # example:
        # row_num multiplied by cell height
        # plus half a cell height (to center the value)

        left_offset = col_num * CELL_DIMS + (CELL_DIMS / 2)
        top_offset = row_num * CELL_DIMS + (CELL_DIMS / 2)
        # print(f"value: {value}, row: {row_num}, top_offset: {top_offset}")
        self.num_text.append(self.canvas.create_text(
            left_offset,
            top_offset,
            text=str(value),
            fill="black",
            font=(f"Helvetica {FONT_SIZE}")
        ))

    def place_guess(self, x, y, value):
        # FYI: this code is awful, though slightly out of scope

        # value dicates position within cell.
        guess_x_offset = ((value - 1) % 3) * (BOARD_DIMS / 9 / 3)
        guess_y_offset = ((value - 1) // 3) * (BOARD_DIMS / 9 / 3)

        top_offset = ((BOARD_DIMS / 9 / 3 / 2)) + x * (BOARD_DIMS / 9) + guess_x_offset
        left_offset = ((BOARD_DIMS / 9 / 3 / 2)) + y * (BOARD_DIMS / 9) + guess_y_offset
        self.guess_text.append(self.canvas.create_text(
            top_offset,
            left_offset,
            text=str(value),
            fill="black",
            font=('Helvetica 15')
        ))

    def _click_canvas(self, event):
        x = event.x
        y = event.y

        for cell in range(len(self.cell_rect)):
            self.canvas.itemconfig(self.cell_rect[cell], fill="white")

        for cell in range(len(self.cell_rect)):
            x0, y0, x1, y1 = self.canvas.coords(self.cell_rect[cell])
            if x > x0 and y > y0 and x < x1 and y < y1: #finding intersect
                self.canvas.itemconfig(self.cell_rect[cell], fill="grey")
                self.selected_cell = cell
                break

    def _select_next_cell(self):
        if self.selected_cell == 80:
            return
        
        for cell in range(len(self.cell_rect)):
            if cell != self.selected_cell:
                continue
            
            self.canvas.itemconfig(self.cell_rect[cell], fill="white")
            self.canvas.itemconfig(self.cell_rect[cell + 1], fill="grey")
            self.selected_cell += 1
            break

    def _key_pressed(self, event):
        if self.selected_cell == None:
            return

        if event.keysym in ["space", "Tab"]:
            self._select_next_cell()
            return
        
        if event.keysym not in ["1", "2", "3", "4", "5", "6", "7", "8", "9"]:
            return
        
        slctd_row = self.selected_cell // 9
        slctd_col = self.selected_cell % 9

        for cell in self.puzzle:
            if cell.location.row == slctd_row and cell.location.column == slctd_col:
                cell.value = event.keysym
        self._update_gui(self.puzzle)
        
    def export(self):
        puzzle_string = "".join([cell.value for cell in self.puzzle])
        print(puzzle_string)

    def play(self):
        self.window.mainloop()

def main():
    puzzle = SudokuPuzzle("457869003009200058000000740000730000620908075000046000061000000730005100500193867")
    game = GUIBoard()
    game.load_puzzle(puzzle.board)
    game.play()

if __name__ == "__main__":
    main()
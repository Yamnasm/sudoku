from tkinter import *

BOARD_DIMS = 600

class SudokuBoard:
    def __init__(self):
        self.window = Tk()
        self.window.title("Sudoku")
        self.window.resizable(False, False)
        self.window.configure(bg="white")
        self.window.bind("<1>", self.click_canvas)
        self.window.bind("<KeyPress>", self.number_on_canvas)

        self.board_container = Frame(
            self.window,
            background="white"
        )

        self.canvas = Canvas(
            self.board_container,
            height=BOARD_DIMS,
            width=BOARD_DIMS
        )

        self.cell_rect = []
        self.create_cells()

        self.box_rect = []
        self.create_boxes()

        self.num_text = []
        self.guess_text = []

        self.selected_cell = None

        self.board_container.pack()
        self.canvas.pack()

    def create_boxes(self):
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
    
    def create_cells(self):
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
    
    def place_number(self, x, y, value):
        top_offset = (BOARD_DIMS / 9 / 2) + x * (BOARD_DIMS / 9)
        left_offset = BOARD_DIMS / 9 / 2 + y * (BOARD_DIMS / 9)
        self.num_text.append(self.canvas.create_text(
            top_offset,
            left_offset,
            text=str(value),
            fill="black",
            font=('Helvetica 40')
        ))

    def place_guess(self, x, y, value):
        # FYI: this code is awful
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

    def click_canvas(self, event):
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

    def number_on_canvas(self, event):
        if self.selected_cell == None:
            return
        x = self.selected_cell % 9
        y = self.selected_cell // 9

        for nums in range(len(self.num_text)):
            x0, y0, x1, y1 = self.canvas.coords(self.cell_rect[self.selected_cell])
            xn, yn = self.canvas.coords(self.num_text[nums])
            if xn > x0 and yn > y0 and xn < x1 and yn < y1: #finding intersect
                self.canvas.delete(self.num_text[nums])
                self.num_text[nums] = None
        self.num_text = [x for x in self.num_text if x is not None] #fixing list

        if event.keycode >= 49 and event.keycode <= 57:
            for nums in range(len(self.guess_text)):
                x0, y0, x1, y1 = self.canvas.coords(self.cell_rect[self.selected_cell])
                xn, yn = self.canvas.coords(self.guess_text[nums])
                if xn > x0 and yn > y0 and xn < x1 and yn < y1: #finding intersect
                    self.canvas.delete(self.guess_text[nums])
                    self.guess_text[nums] = None
            self.guess_text = [x for x in self.guess_text if x is not None] #fixing list
        
        if event.keycode >= 97 and event.keycode <= 105:
            self.place_guess(x, y, int(event.char))
        if event.keycode >= 49 and event.keycode <= 57:
            self.place_number(x, y, int(event.char))
        
    def play(self):
        self.window.mainloop()

def main():
    game = SudokuBoard()
    game.play()

if __name__ == "__main__":
    main()
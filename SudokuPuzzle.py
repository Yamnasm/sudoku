from dataclasses import dataclass

@dataclass
class Position:
    row:    int
    column: int
    box:    int

@dataclass
class Cell:
    value:    str
    guesses:  str
    location: Position

class SudokuPuzzle:
    def __init__(self, puz_string = None) -> None:
        self.board = []
        if puz_string == None:
            puz_string = "0"*81

        for index in range(81):

            row = index // 9
            col = index % 9
            box = ((row // 3) * 3) + ((col // 3) % 3)

            self.board.append(Cell(
                value    = puz_string[index],
                guesses  = [],
                location = Position(
                    row    = row,
                    column = col,
                    box    = box
                )
            ))
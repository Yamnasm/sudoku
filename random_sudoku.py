# extracts the sudoku game from https://www.websudoku.com/

import argparse
import requests
from bs4 import BeautifulSoup

parser = argparse.ArgumentParser()
parser.add_argument("-b", "--batch", action="store_true")
args = parser.parse_args()

#sudoku difficulty levels:
EASY   = "1"
MEDIUM = "2"
HARD   = "3"
EVIL   = "4"

def retrieve_page(difficulty = "1"):
    sudoku_url = requests.get(f"https://west.websudoku.com/?level={difficulty}")
    html_soup = BeautifulSoup(sudoku_url.text, "html.parser")
    return html_soup

def parse_table(html_page):
    raw_table = html_page.find("table", id="puzzle_grid")
    sudoku_table = raw_table.find_all("td")
    sudoku_str_arr = []

    # the single most painful find_all loop I've ever had to write.
    for table_element in sudoku_table:
        element = table_element.find_all("input")[0]
        if element.has_attr("value"):
            sudoku_str_arr.append(element["value"])
        else:
            sudoku_str_arr.append("0")
    return "".join(sudoku_str_arr)

def get_random_sudoku_puzzle():
    html_page = retrieve_page(EVIL)
    sudoku_string = parse_table(html_page)
    return sudoku_string

def get_puzzle_batch(loop = 50):
    puzzle_list = []
    for _ in range(loop):
        try:
            puzzle_list.append(get_random_sudoku_puzzle())
        except requests.exceptions.SSLError:
            break
    with open("test_puzzles.txt", 'w') as file:
        for puz in puzzle_list:
            file.write(f"{puz}\n")

def main():
    if args.batch:
        get_puzzle_batch()
        return
    print(get_random_sudoku_puzzle())

if __name__ == "__main__":
    main()
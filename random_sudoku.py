# extracts the sudoku game from https://www.websudoku.com/

import requests
from bs4 import BeautifulSoup

def retrieve_page():
    sudoku_url = requests.get("https://west.websudoku.com/?")
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
    html_page = retrieve_page()
    sudoku_string = parse_table(html_page)
    return sudoku_string

if __name__ == "__main__":
    get_random_sudoku_puzzle()
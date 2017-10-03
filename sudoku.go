package main

import(
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	fmt.Println("\033[H\033[2J")
	sudoku := Sudoku{}
	//sudoku = populateRandomly(sudoku)
	sudoku = populateManualy(sudoku)
	printSudoku(sudoku)
	if(isSudokuValid(sudoku)) {
		fmt.Println("\nThis sudoku is valid\n")
	}else {
		fmt.Println("\nThis sudoku is not valid\n")
	}
}

type Sudoku struct {
	grid [9][9] int
}

//-------MANUALY

func populateManualy(sudoku Sudoku) Sudoku {
	sudoku.grid =  [9][9] int{
		{5,3,0/*|*/,0,7,0/*|*/,0,0,0},
		{6,0,0/*|*/,1,9,5/*|*/,0,0,0},
		{4,9,8/*|*/,0,0,0/*|*/,0,6,0},
		//---------------------------
		{8,0,0/*|*/,0,6,0/*|*/,0,0,3},
		{0,0,0/*|*/,8,0,3/*|*/,0,0,1},
		{7,0,0/*|*/,0,2,0/*|*/,0,0,6},
		//---------------------------
		{0,6,0/*|*/,0,0,0/*|*/,2,8,0},
		{0,0,0/*|*/,4,1,9/*|*/,0,0,5},
		{0,0,0/*|*/,0,8,0/*|*/,0,7,9}}
	return sudoku
}

//-------RANDOM

func populateRandomly(sudoku Sudoku) Sudoku {
	for indexX, valueX := range sudoku.grid {
		for indexY, _ := range valueX {
			sudoku.grid[indexX][indexY] = random(1, 9)
		}
	}
	return sudoku
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min) + min
}

//--------CHECK

func isSudokuValid(sudoku Sudoku) bool {
	return checkHorizontaly(sudoku) && checkVerticaly(sudoku) && checkCases(sudoku)
}

func checkHorizontaly(sudoku Sudoku) bool {
	values := [9] int{0,0,0,0,0,0,0,0,0}
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if sudoku.grid[y][x] != 0 {
				for _, value := range values {
					if value == sudoku.grid[y][x] {
						return false
					} 
				}
				values[x] = sudoku.grid[y][x]
			}
		}
		values = [9] int{0,0,0,0,0,0,0,0,0}
	}
	return true
}

func checkVerticaly(sudoku Sudoku) bool {
	values := [9] int{0,0,0,0,0,0,0,0,0}
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if sudoku.grid[y][x] != 0 {
				for _, value := range values {
					if value == sudoku.grid[y][x] {
						return false
					} 
				}
				values[y] = sudoku.grid[y][x]
			}
		}
		values = [9] int{0,0,0,0,0,0,0,0,0}
	}
	return true
}

func checkCases(sudoku Sudoku) bool {
	values := [9] int{0,0,0,0,0,0,0,0,0}
	for y := 0; y < 9; y = y+3 {
		for x := 0; x < 9; x = x+3 {
			for i := 0; i < 3; i++ {
				for j := 0; x < 3; x++ {
					if sudoku.grid[y+i][x+j] != 0 {
						for _, value := range values {
							if value == sudoku.grid[y+i][x+j] {
								return false
							} 
						}
						values[(i+1)*(j+1)] = sudoku.grid[y+i][x+j]
					}
				}
			}
			values = [9] int{0,0,0,0,0,0,0,0,0}	
		}
	}
	return true
}

//--------PRINT

func printSudoku(sudoku Sudoku) {
	line := ""
	for indexX, valueX := range sudoku.grid {
		if(indexX == 3 || indexX == 6){
			for i := 0; i < 14; i++ {
    			line += "--"
			}
			fmt.Println(line)
			line = ""
		}
		for indexY, valueY := range valueX {
			if(indexY == 3 || indexY == 6) {
				line += "|"
			}
			line += " " + strconv.Itoa(valueY) + " "
		}
		fmt.Println(line)
		line = ""
	}
}
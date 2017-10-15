package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"os"
)

func main() {
	fmt.Println("\033[H\033[2J") // Clear
	sudoku := Sudoku{}

	//sudoku = populateRandomly(sudoku) // Randomly generate the sudoku
	populateManualy(&sudoku) // You can change it in the func populateManualy
	sudokuBefore := sudoku
	if !isSudokuValid(&sudoku) {
		printSudoku(&sudoku, false)
		fmt.Println("\033[31mThe sudoku is not valid\033[0m")
	} else {
		var coord = getChangeableCoordinates(&sudoku) // Search all 0 in the grid

		start := time.Now()
		//solveSudoku(Sudoku,
		//			 []Coord,
		//			 position #must be 0#, back #must be false#,
		//			 Sudoku #it'll not change#,
		//			 timeLaps #time in Millisecond or -1 for no timelaps#)
		timeLaps := -1

		if(len(os.Args) != 1){
			argTL, err := strconv.Atoi(os.Args[1])			
			if(err == nil){			
				timeLaps = argTL
			}
		}

		position := 0
		back := false
		solveSodoku(&sudoku, coord, &position, &back, &sudokuBefore, &timeLaps)
		elapsed := time.Since(start) // Time spend by solveSudoku

		fmt.Println("The solver function took\033[31m", elapsed)
		fmt.Println("\033[0m")
	}
}

type Sudoku struct {
	grid [9][9]int
}

type Coord struct {
	y int
	x int
}

//-------MANUALY

func populateManualy(sudoku *Sudoku) {
	// Good easy
	sudoku.grid = [9][9]int{
		{5, 3, 0 /*|*/, 0, 7, 0 /*|*/, 0, 0, 0},
		{6, 0, 0 /*|*/, 1, 9, 5 /*|*/, 0, 0, 0},
		{0, 9, 8 /*|*/, 0, 0, 0 /*|*/, 0, 6, 0},
		//---------------------------
		{8, 0, 0 /*|*/, 0, 6, 0 /*|*/, 0, 0, 3},
		{4, 0, 0 /*|*/, 8, 0, 3 /*|*/, 0, 0, 1},
		{7, 0, 0 /*|*/, 0, 2, 0 /*|*/, 0, 0, 6},
		//---------------------------
		{0, 6, 0 /*|*/, 0, 0, 0 /*|*/, 2, 8, 0},
		{0, 0, 0 /*|*/, 4, 1, 9 /*|*/, 0, 0, 5},
		{0, 0, 0 /*|*/, 0, 8, 0 /*|*/, 0, 7, 9}}

	// DEAMON
	//sudoku.grid = [9][9] int{
	//	{0,0,0/*|*/,0,0,9/*|*/,1,0,0},
	//	{0,0,9/*|*/,0,3,0/*|*/,5,0,0},
	//	{0,1,0/*|*/,0,0,6/*|*/,0,7,9},
	//	//---------------------------
	//	{0,0,0/*|*/,0,6,0/*|*/,3,0,8},
	//	{0,9,0/*|*/,3,0,7/*|*/,0,1,0},
	//	{1,0,6/*|*/,0,4,0/*|*/,0,0,0},
	//	//---------------------------
	//	{7,6,0/*|*/,4,0,0/*|*/,0,2,0},
	//	{0,0,8/*|*/,0,1,0/*|*/,9,0,0},
	//	{0,0,2/*|*/,7,0,0/*|*/,0,0,0}}

	// Not good
	//sudoku.grid =  [9][9] int{
	//	{5,3,0/*|*/,0,7,0/*|*/,0,0,7},
	//	{6,0,0/*|*/,1,9,5/*|*/,0,0,0},
	//	{0,9,8/*|*/,0,0,0/*|*/,0,6,0},
	//	//---------------------------
	//	{8,0,0/*|*/,0,6,0/*|*/,0,0,3},
	//	{4,0,0/*|*/,8,0,3/*|*/,0,0,1},
	//	{7,0,0/*|*/,0,2,0/*|*/,0,0,6},
	//	//---------------------------
	//	{0,6,0/*|*/,0,0,0/*|*/,2,8,0},
	//	{0,0,0/*|*/,4,1,9/*|*/,0,0,5},
	//	{0,0,0/*|*/,0,8,0/*|*/,0,7,9}}

	// Impossible
	//sudoku.grid =  [9][9] int{
	//	{5,3,0/*|*/,0,7,0/*|*/,0,0,0},
	//	{6,0,0/*|*/,1,9,5/*|*/,0,0,0},
	//	{4,9,8/*|*/,0,0,0/*|*/,0,6,0},
	//	//---------------------------
	//	{8,0,0/*|*/,0,6,0/*|*/,0,0,3},
	//	{0,0,0/*|*/,8,0,3/*|*/,0,0,1},
	//	{7,0,0/*|*/,0,2,0/*|*/,0,0,6},
	//	//---------------------------
	//	{0,6,0/*|*/,0,0,0/*|*/,0,8,0},
	//	{0,0,0/*|*/,4,1,9/*|*/,0,0,5},
	//	{0,0,0/*|*/,0,8,0/*|*/,0,7,9}}
}

//-------RANDOM

func populateRandomly(sudoku *Sudoku) {
	for indexX, valueX := range sudoku.grid {
		for indexY, _ := range valueX {
			if random(1, 10) == 4 {
				sudoku.grid[indexX][indexY] = random(1, 9)
			}
		}
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//--------CHECK

func isSudokuValid(sudoku *Sudoku) bool {
	//return checkHorizontaly(sudoku) && checkVerticaly(sudoku) && checkCases(sudoku) // Check the 3 possibilities
	if !checkHorizontaly(sudoku) {
		return false
	}
	if !checkVerticaly(sudoku) {
		return false
	}
	if !checkCases(sudoku) {
		return false
	}
	return true
}

func checkHorizontaly(sudoku *Sudoku) bool {
	/*
		For on y axe then x axe and check if there is 2 values
		If it is then return false
		If it is not continue on the next y axe
	*/
	var values []int
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if sudoku.grid[y][x] != 0 {
				for _, value := range values {
					if value == sudoku.grid[y][x] {
						return false
					}
				}
				values = append(values, sudoku.grid[y][x])
			}
		}
		values = nil
	}
	return true
}

func checkVerticaly(sudoku *Sudoku) bool {
	/*
		For on x axe then y axe and check if there is 2 same values
		If it is then return false
		If it is not continue on the next x axe
	*/
	var values []int
	for x := 0; x < 9; x++ {
		for y := 0; y < 9; y++ {
			if sudoku.grid[y][x] != 0 {
				for _, value := range values {
					if value == sudoku.grid[y][x] {
						return false
					}
				}
				values = append(values, sudoku.grid[y][x])
			}
		}
		values = nil
	}
	return true
}

func checkCases(sudoku *Sudoku) bool {
	/*
		Check 3 by 3

		1OO|2OO|3OO
		OOO|OOO|OOO
		OOO|OOO|OOO
		-----------
		4OO|5OO|6OO
		OOO|OOO|OOO
		OOO|OOO|OOO
		-----------
		7OO|8OO|9OO
		OOO|OOO|OOO
		OOO|OOO|OOO

		1-9 is where we start
		Then we check all the case in this order

		123
		456
		789

		Then check if there no 2 same values
		If it is then return false
		If it is not continue on the next case
	*/
	var values []int
	for y := 0; y < 9; y = y + 3 {
		for x := 0; x < 9; x = x + 3 {
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if sudoku.grid[y+i][x+j] != 0 {
						for _, value := range values {
							if value == sudoku.grid[y+i][x+j] {
								return false
							}
						}
						values = append(values, sudoku.grid[y+i][x+j])
					}
				}
			}
			values = nil
		}
	}
	return true
}

//-------SOLVE

func solveSodoku(sudoku *Sudoku, coord []Coord, position *int, back *bool, sudokuBefore *Sudoku, timeLaps *int) {
	var wasBack = false // Use to print the sudoku
	// Solver
	if *position == -1 { // Impossible sudoku
		printSudoku(sudokuBefore, false)
		fmt.Println("\n\033[31mThere is no solution for this sudoku\033[0m\n")
		return
	}
	if sudoku.grid[coord[*position].y][coord[*position].x] == 9 && *back == true { // If the lastest action was back and the value is 9
		sudoku.grid[coord[*position].y][coord[*position].x] = 0
		*back = false
		wasBack = true
		*position--
	} else if isSudokuValid(sudoku) {
		var start = sudoku.grid[coord[*position].y][coord[*position].x] // This is the value of the cell
		var check = 0                                                 // This is the value of the incrementation

		for i := start + 1; i <= 9; i++ { // Start at the value of the cell +1
			check = i
			sudoku.grid[coord[*position].y][coord[*position].x] = i // The cell take the value of the incrementation
			if isSudokuValid(sudoku) {                            // Means that the sodoku is good : let's take another one
				if *position == len(coord)-1 {
					// YE4H !!! It's done
					fmt.Println("\033[H\033[2J")
					fmt.Println("\n\033[33mBefore : \033[0m\n")
					printSudoku(sudokuBefore, false)
					fmt.Println("\n\033[32mSolution : \033[0m\n")
					printSudoku(sudoku, false)
					return // Stop
				}
				break // Stop the for because the sudoku is valid
			}
		}
		if check == 9 && !isSudokuValid(sudoku) { // Means that the sudoku is not good : let's go back
			*back = true
			sudoku.grid[coord[*position].y][coord[*position].x] = 0
			*position--
		} else { // Good
			*position++
		}
	}
	// TimeLaps
	if *timeLaps != -1 {
		fmt.Println("\033[H\033[2J") // Clear
		if *back == true {
			printSudoku(sudoku, coord[*position])
		} else {
			if wasBack == true {
				printSudoku(sudoku, coord[*position])
			} else {
				printSudoku(sudoku, coord[*position-1])
			}
		}
		time.Sleep(time.Duration(*timeLaps) * time.Millisecond)
	}
	// Recursiv
	solveSodoku(sudoku, coord, position, back, sudokuBefore, timeLaps)
}

//-------COORD

func getChangeableCoordinates(sudoku *Sudoku) []Coord {
	var coord []Coord
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if sudoku.grid[y][x] == 0 {
				coord = append(coord, Coord{y, x})
			}
		}
	}
	return coord
}

//--------PRINT

func printSudoku(sudoku *Sudoku, coord interface{}) {
	line := ""
	for indexY, valueY := range sudoku.grid {
		if indexY == 0 {
			line += "╔═══╦═══╦═══╦═══╦═══╦═══╦═══╦═══╦═══╗"
			fmt.Println(line)
			line = ""
		}

		if indexY == 3 || indexY == 6 {
			line += "╠═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╬═══╣"
			fmt.Println(line)
			line = ""
		}
		for indexX, valueX := range valueY {
			if indexX == 0 {
				line += "║"
			}
			if indexX == 3 || indexX == 6 {
				line += "║"
			}

			if valueX != 0 {
				line += " "
				switch coord.(type) {
                    case Coord:
                    	coord2 := coord.(Coord)
                    	if(coord2.x == indexX && coord2.y == indexY){
                     	  	line += "\033[31m" + strconv.Itoa(valueX) + "\033[0m"
                    	}else {
                    		line += strconv.Itoa(valueX)	
                    	}
                    default:
                        line += strconv.Itoa(valueX)
            	}
				line += " "
			} else {
				line += " ░ "
			}

			if indexX != 2 && indexX != 5 && indexX != 8 {
				line += "│"
			}
		}
		line += "║"
		fmt.Println(line)
		line = ""
		if indexY != 2 && indexY != 5 && indexY != 8 {
			line += "╠───┼───┼───╬───┼───┼───╬───┼───┼───╣"
			fmt.Println(line)
			line = ""
		}
	}
	line += "╚═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╝"
	fmt.Println(line)
}

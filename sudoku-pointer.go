package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"os"
	"io/ioutil"
)

var timeLaps = -1
var back = false
var coord []Coord
var listSudoku []Sudoku
var sudoku Sudoku
const path = "./example/"

func main() {
	fmt.Println("\033[H\033[2J") // Clear
	if (len(os.Args) == 2){
		if(os.Args[1] == "help"){
			fmt.Println("\033[31mUsage\033[0m : go run \033[32mfile\033[0m \033[32mtimeLaps\033[0m \033[32mtype\033[0m\n")
			fmt.Println("\033[32mfile\033[0m is the file you selected (\033[33msudoku.go\033[0m or \033[33msudoku-pointer.go\033[0m)")
			fmt.Println("\033[32mtimeLaps\033[0m is the intervall in milliSecond (\033[33m-1\033[0m is default and for no timeLaps)")
			return
		}
	}

	sudoku = Sudoku{}
	populateManualy() // You can change it in the func populateManualy

	for _, sudoo := range listSudoku {
		sudoku = sudoo
		doWhatYouWantToDo()
		//printSudoku(&sudoku, false)
	}
}

func doWhatYouWantToDo() {
	getChangeableCoordinates() // Search all 0 in the grid

	if(len(os.Args) != 1){
		argTL, err := strconv.Atoi(os.Args[1])
		if(err == nil){
			timeLaps = argTL
		}
	}

	start := time.Now()
	solveSodoku(0)
	elapsed := time.Since(start) // Time spend by solveSudoku

	fmt.Println("The solver function took\033[31m", elapsed)
	fmt.Println("\033[0m")
}

type Sudoku struct {
	grid [9][9]int
}

type Coord struct {
	y int
	x int
}

//-------MANUALY

func populateManualy() {
	files, err := ioutil.ReadDir(path)
    if err != nil {
		fmt.Println("err")
        os.Exit(0)
    }

    for _, f := range files {
        content, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			fmt.Println("err")
			os.Exit(0)
		}

		sudoku = Sudoku{}
		x := 0
		y := 0

		for _, ch := range content{
			str := string(ch)
			if(str == "."){
				sudoku.grid[y][x] = 0
			}else if str != "\n" {
				sudoku.grid[y][x], err = strconv.Atoi(str)
				if err != nil {
					fmt.Println("err")
				}
			}
			if(x > 8){
				x = 0
				y = y+1
			}else {
				x = x+1
			}
		}
		listSudoku = append(listSudoku, sudoku)
	}
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

func isSudokuValid() bool {
	if !checkHorizontaly() {
		return false
	}
	if !checkVerticaly() {
		return false
	}
	if !checkCases() {
		return false
	}
	return true
}

func checkHorizontaly() bool {
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

func checkVerticaly() bool {
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

func checkCases() bool {
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

func solveSodoku(position int) {
	if position == -1 { // Impossible sudoku
		fmt.Println("\n\033[31mThere is no solution for this sudoku\033[0m\n")
		return
	}
	if sudoku.grid[coord[position].y][coord[position].x] == 9 && back == true { // If the lastest action was back and the value is 9
		sudoku.grid[coord[position].y][coord[position].x] = 0
		back = false
		position--
	} else if isSudokuValid() {
		var start = sudoku.grid[coord[position].y][coord[position].x] // This is the value of the cell
		var check = 0                                                 // This is the value of the incrementation

		for i := start + 1; i <= 9; i++ { // Start at the value of the cell +1
			check = i
			sudoku.grid[coord[position].y][coord[position].x] = i // The cell take the value of the incrementation
			if isSudokuValid() {                            // Means that the sodoku is good : let's take another one
				if position == len(coord)-1 {
					// YE4H !!! It's done
					fmt.Println("\n\033[32mSolution : \033[0m\n")
					printSudoku(&sudoku, false)
					return // Stop
				}
				break // Stop the for because the sudoku is valid
			}
		}
		if check == 9 && !isSudokuValid() { // Means that the sudoku is not good : let's go back
			back = true
			sudoku.grid[coord[position].y][coord[position].x] = 0
			position--
		} else { // Good
			position++
		}
	}
	// Recursiv
	solveSodoku(position)
}

//-------COORD

func getChangeableCoordinates() {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if sudoku.grid[y][x] == 0 {
				coord = append(coord, Coord{y, x})
			}
		}
	}
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
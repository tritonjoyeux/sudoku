package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Sudoku struct {
	Grid    [9][9]int
	Name    string
	Elapsed time.Duration
}

func displaySudoku(ch2 chan Sudoku, wg *sync.WaitGroup) {
	for {
		var base Sudoku
		base = <-ch2

		fmt.Println(base.Name)

		line := ""
		for indexY, valueY := range base.Grid {
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
					line += " " + strconv.Itoa(valueX) + " "
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
		fmt.Println("The solver function took\033[31m", base.Elapsed)
		fmt.Println("\033[0m")
		wg.Done()
	}
}

func (this *Sudoku) Check() bool {
	for x := 0; x < 9; x++ {
		acc := make(map[int]bool)
		for y := 0; y < 9; y++ {
			val := this.Grid[x][y]
			if acc[val] && val != 0 {
				return false
			}
			acc[val] = true
		}
	}

	for y := 0; y < 9; y++ {
		acc := make(map[int]bool)
		for x := 0; x < 9; x++ {
			val := this.Grid[x][y]
			if acc[val] && val != 0 {
				return false
			}
			acc[val] = true
		}
	}

	for cadX := 0; cadX < 3; cadX++ {
		for cadY := 0; cadY < 3; cadY++ {
			acc := make(map[int]bool)
			for x := cadX * 3; x < 3; x++ {
				for y := cadY * 3; y < 3; y++ {
					val := this.Grid[x][y]
					if acc[val] && val != 0 {
						return false
					}
					acc[val] = true
				}
			}
		}
	}
	return true
}

func (this *Sudoku) Solve() {
	coords := this.getMissingsNumbers()
	this.solveRecursion(coords)
}

func (this *Sudoku) getMissingsNumbers() (res [][2]int) {
	for ky, vy := range this.Grid {
		for kx, vx := range vy {
			if vx == 0 {
				add := [2]int{ky, kx}
				res = append(res, add)
			}
		}
	}
	return
}

func (this *Sudoku) solveRecursion(coords [][2]int) bool {
	if len(coords) == 0 {
		return true
	}
	y := coords[0][0]
	x := coords[0][1]
	for n := 1; n <= 9; n++ {
		if this.checkCoord(y, x, n) {
			this.Grid[y][x] = n
			if this.solveRecursion(coords[1:]) {
				return true
			}
			this.Grid[y][x] = 0
		}
	}
	return false
}

func (this *Sudoku) checkCoord(cy int, cx int, nVal int) bool {
	// Line
	for x := 0; x < 9; x++ {
		val := this.Grid[cy][x]
		if val == nVal {
			return false
		}
	}

	// Col
	for y := 0; y < 9; y++ {
		val := this.Grid[y][cx]
		if val == nVal {
			return false
		}
	}

	// square
	by := cy - (cy % 3)
	bx := cx - (cx % 3)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			val := this.Grid[by+y][bx+x]
			if val == nVal {
				return false
			}
		}
	}
	return true
}

func getFolderSudoku() []Sudoku {
	var listSudoku []Sudoku
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	sudoku := Sudoku{}
	for _, f := range files {
		content, err := ioutil.ReadFile(path + f.Name())
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		sudoku = Sudoku{}
		sudoku.Name = f.Name()
		x := 0
		y := 0

		for _, ch := range content {
			str := string(ch)
			if str == "." {
				sudoku.Grid[y][x] = 0
			} else if str != "\n" {
				sudoku.Grid[y][x], _ = strconv.Atoi(str)
			}
			if x > 8 {
				x = 0
				y = y + 1
			} else {
				x = x + 1
			}
		}
		listSudoku = append(listSudoku, sudoku)
	}
	return listSudoku
}

const path = "./example/"

func main() {
	var wg sync.WaitGroup

	listSudoku := getFolderSudoku()

	var ch chan Sudoku
	ch = make(chan Sudoku)

	var ch2 chan Sudoku
	ch2 = make(chan Sudoku)

	for i := 0; i < runtime.NumCPU(); i++ {
		go threadSudoku(ch, &wg, ch2)
	}

	go displaySudoku(ch2, &wg)

	for _, base := range listSudoku {
		wg.Add(1)
		ch <- base
	}
	wg.Wait()
}

func threadSudoku(ch chan Sudoku, wg *sync.WaitGroup, ch2 chan Sudoku) {
	for {
		var base Sudoku
		base = <-ch
		start := time.Now()
		base.Solve()
		base.Elapsed = time.Since(start) // Time spend by solveSudoku
		ch2 <- base
	}
}

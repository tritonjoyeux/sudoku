package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
	"runtime/pprof"
	"log"
)

type Sudoku struct {
	Grid [9][9]int
	Name string
}

func displaySudoku(ch2 chan Sudoku, wg *sync.WaitGroup) {
	for {
		var base Sudoku
		base = <-ch2
		file, err := os.Create("./results/" + base.Name + ".txt")
		if(err != nil) {
			log.Fatal("could not create the result file " + base.Name + ".txt")
		}
		defer file.Close()
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
					fmt.Fprintf(file, strconv.Itoa(valueX))
				} else {
					line += " ░ "
				}

				if indexX != 2 && indexX != 5 && indexX != 8 {
					line += "│"
				}
			}
			line += "║"
			fmt.Println(line)
			fmt.Fprintf(file, "\n")
			line = ""
			if indexY != 2 && indexY != 5 && indexY != 8 {
				line += "╠───┼───┼───╬───┼───┼───╬───┼───┼───╣"
				fmt.Println(line)
				line = ""
			}
		}
		line += "╚═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╩═══╝"
		fmt.Println(line)
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

func getFolderSudoku(ch chan Sudoku, wg *sync.WaitGroup, ch2 chan Sudoku, chStart chan string) {
	for {
		var pathGlob string
		pathGlob = <-chStart

		files, err := ioutil.ReadDir(pathGlob)
		if err != nil {
			log.Fatal("could not read Dir " + pathGlob + ": ", err)
		}
		sudoku := Sudoku{}
		for _, f := range files {
			content, err := ioutil.ReadFile(pathGlob + f.Name())
			if err != nil {
				log.Fatal("could not read File " + f.Name() + ": ", err)
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
			wg.Add(1)
			ch <- sudoku
		}
		wg.Done()
	}
}

const path = "./example/"

func main() {
	start := time.Now()
	f, err := os.Create("perf_cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	var wg sync.WaitGroup

	var chStart chan string
	chStart = make(chan string)

	var ch chan Sudoku
	ch = make(chan Sudoku)

	var ch2 chan Sudoku
	ch2 = make(chan Sudoku)

	go getFolderSudoku(ch, &wg, ch2, chStart)

	for i := 0; i < runtime.NumCPU(); i++ {
		go threadSudoku(ch, &wg, ch2)
	}

	go displaySudoku(ch2, &wg)

	wg.Add(1)
	chStart <- path
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("The solver function took\033[31m", elapsed)
	fmt.Println("\033[0m")

	f, err = os.Create("perf_mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC()
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}

func threadSudoku(ch chan Sudoku, wg *sync.WaitGroup, ch2 chan Sudoku) {
	for {
		var base Sudoku
		base = <-ch
		base.Solve()
		ch2 <- base
	}
}

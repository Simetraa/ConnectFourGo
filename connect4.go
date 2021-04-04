package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	empty  = iota
	yellow = iota
	red    = iota
)

const (
	ansiYellow     = "\x1b[93m"
	ansiRed        = "\x1b[31m"
	ansiBlue       = "\x1b[96m"
	ansiBold       = "\x1b[1m"
	ansiUnderlined = "\x1b[4m"
	ansiReset      = "\x1b[0m"
)

const (
	rowCount = 6
	colCount = 7
)

type Board struct {
	board         [][]int
	currentPlayer int
}

func (b *Board) init() {
	matrix := make([][]int, rowCount)
	for i := 0; i < len(matrix); i++ {
		matrix[i] = make([]int, colCount)
	}
	b.currentPlayer = yellow // set starting player
	b.board = matrix
}

func (b *Board) getBoard() string {
	var output string

	output += ansiBold + ansiUnderlined + ansiBlue
	var header []string
	for i := 0; i < colCount; i++ {
		header = append(header, strconv.Itoa(i))
	}
	output += strings.Join(header, " ")
	output += fmt.Sprintf("%s\n", ansiReset)
	for _, row := range b.board {
		for _, counter := range row {
			switch counter {
			case yellow:
				output += ansiYellow + "● "
			case red:
				output += ansiRed + "● "
			case empty:
				output += ansiReset + "◯ "
			}
		}
		output += "\n"
	}
	return output
}

func (b *Board) dropCounter(colour, x int) (int, error) {
	if x < 0 || x > colCount {
		return -1, errors.New("column is out of range")
	}
	for y := len(b.board) - 1; y >= 0; y-- {
		if b.board[y][x] == empty {
			b.board[y][x] = colour
			return y, nil
		}

	}
	return -1, errors.New("column is full")
}

func (b *Board) cellExists(x, y int) bool {
	return ((x >= 0 && x <= colCount-1) && (y >= 0 && y <= rowCount-1))
}

func (b *Board) getLength() int {
	var count int
	for _, row := range b.board {
		for _, counter := range row {
			if counter != empty {
				count++
			}
		}
	}
	return count

}

func (b *Board) togglePlayer() {
	switch b.currentPlayer {
	case yellow:
		b.currentPlayer = red
	case red:
		b.currentPlayer = yellow
	}
}

func (b *Board) isWin(colour, xPos, yPos int) bool {
	count := 0
	for x, y := -3, -3; x <= 3; x, y = x+1, y+1 {
		if !b.cellExists(xPos+x, yPos+y) {
			count = 0
			continue
		}
		if b.board[yPos+y][xPos+x] == colour {
			count++
		} else { // reset counter
			count = 0
		}
		if count == 4 {
			return true
		}
	}
	count = 0
	for x, y := -3, 3; x <= 3; x, y = x+1, y-1 {
		if !b.cellExists(xPos+x, yPos+y) {
			count = 0
			continue
		}
		if b.board[yPos+y][xPos+x] == colour {
			count++
		} else {
			count = 0
		}
		if count == 4 {
			return true
		}
	}
	count = 0
	for x := -3; x <= 3; x++ {
		if !b.cellExists(xPos+x, yPos) {
			count = 0
			continue
		}
		if b.board[yPos][xPos+x] == colour {
			count++
		} else {
			count = 0
		}
		if count == 4 {
			return true
		}
	}

	count = 0
	for y := -3; y <= 3; y++ {
		if !b.cellExists(xPos, yPos+y) {
			count = 0
			continue
		}
		if b.board[yPos+y][xPos] == colour {
			count++
		} else {
			count = 0
		}
		if count == 4 {
			return true
		}
	}
	return false
}

func main() {
	z := Board{}
	z.init()
	for {
		fmt.Println(z.getBoard())
		reader := bufio.NewReader(os.Stdin)
		switch z.currentPlayer {
		case yellow:
			fmt.Printf("%s⚫Enter column: %s", ansiYellow, ansiReset)
		case red:
			fmt.Printf("%s⚫Enter column: %s", ansiRed, ansiReset)
		}
		xInput, _ := reader.ReadString('\n')
		xInput = strings.TrimSuffix(xInput, "\n")
		xInt, err := strconv.Atoi(xInput)
		if err != nil {
			fmt.Println("Invalid entry, please try again")
			continue
		}
		restingY, err := z.dropCounter(z.currentPlayer, xInt)
		if err != nil {
			fmt.Println("Could not drop counter there,", err)
			continue
		}
		if z.getLength() == rowCount*colCount {
			// If all positions in the board are filled.
			fmt.Printf("%sDraw!%s\n", ansiBlue, ansiReset)
			break
		}
		if z.isWin(z.currentPlayer, xInt, restingY) {
			fmt.Println(z.getBoard())
			switch z.currentPlayer {
			case yellow:
				fmt.Printf("%sYellow won!%s\n", ansiYellow, ansiReset)
			case red:
				fmt.Printf("%sRed won!%s\n", ansiRed, ansiReset)
			}
			break
		} else {
			z.togglePlayer()
		}
	}
}

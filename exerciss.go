package main

import (
	"exerciss/algorithms"
	"fmt"
	"time"
)

func main() {
	board := [4][4]int{} //board[row][col]
	board[0][3] = 3
	board[1][0] = 4
	board[2][1] = 1
	board[2][2] = 3
	board[3][0] = 3
	board[3][2] = 2
	//fmt.Println(board)
	before := time.Now()
	answ := algorithms.SolveAlgX(board)
	after := time.Now()
	dur := after.Sub(before)
	fmt.Println(answ)
	fmt.Printf("duration algorithmX %d us\n", dur.Microseconds())
	//SolveAlgX(board)
}

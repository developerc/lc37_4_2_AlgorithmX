package main

import (
	"exerciss/algorithms"
	"fmt"
	"time"
)

func main() {
	//-------- AlgorithmX ---------
	board := [4][4]int{} //board[row][col]
	board[0][3] = 3
	board[1][0] = 4
	board[2][1] = 1
	board[2][2] = 3
	board[3][0] = 3
	board[3][2] = 2
	//fmt.Println(board)
	before := time.Now()
	answ := algorithms.SolveAlgX4_3(board)
	after := time.Now()
	dur := after.Sub(before)
	fmt.Println(answ)
	fmt.Printf("duration algorithmX %d ns\n", dur.Nanoseconds())
	//-------- DFS ---------
	board[0][3] = 3
	board[1][0] = 4
	board[2][1] = 1
	board[2][2] = 3
	board[3][0] = 3
	board[3][2] = 2
	before = time.Now()
	answ = algorithms.SolveDfs(board)
	after = time.Now()
	dur = after.Sub(before)
	fmt.Println(answ)
	fmt.Printf("duration DFS %d ns\n", dur.Nanoseconds())
}

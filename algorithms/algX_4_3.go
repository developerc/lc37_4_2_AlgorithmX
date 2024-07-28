package algorithms

import (
	"fmt"
	"os"
	"strconv"
)

var iter int

// способ решения побитовыми операциями
func SolveAlgX4_3(board [4][4]int) [4][4]int {
	var colRowTable [64]uint64
	fillTable(&board, &colRowTable)
	//colRowTable := fillTable(&board)
	printTable("net"+strconv.Itoa(n)+".txt", colRowTable)
	return board
}

func fillTable(board *[4][4]int, colRowTable *[64]uint64) /*[64]uint64*/ {
	//var colRowTable [64]uint64
	var mul1 uint64 = 1
	mul2 := [4]uint64{1 << 16, 1 << 20, 1 << 24, 1 << 28}

	//mul1 = 1
	for j := 0; j < 16; j++ {
		for i := 0; i < 4; i++ {
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul1    //ограничения в ячейках
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul2[i] //ограничения в строках
		}
		mul1 = mul1 << 1
		if j == 3 || j == 7 || j == 11 {
			for k := 0; k < 4; k++ {
				mul2[k] = mul2[k] << 1
			}
		}
	}

	//return colRowTable
}

func printTable(fileName string, colRowTable [64]uint64) {
	fo, err := os.Create(fileName) // open output file
	if err != nil {
		panic(err)
	}
	defer fo.Close() // close fo on exit and check for its returned error
	for i := 0; i < 64; i++ {
		fo.WriteString(fmt.Sprintf("%064b\n", colRowTable[i]))
	}
	/*var x uint64 = 1 << 12
	x = x >> 10
	fmt.Println(bits.OnesCount(uint(x)))
	//fmt.Printf("%064b", x)
	//fo.WriteString(fmt.Printf("%064b", x))
	fo.WriteString(fmt.Sprintf("%064b", x))*/
}

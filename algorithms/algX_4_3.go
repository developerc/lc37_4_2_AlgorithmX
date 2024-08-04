package algorithms

import (
	"fmt"
	"os"
	"strconv"
)

var iter int
var Nex int

// способ решения побитовыми операциями
func SolveAlgX4_3(board [4][4]int) [4][4]int {
	var colRowTable [64]uint64
	fillTable(&board, &colRowTable)
	//colRowTable := fillTable(&board)
	printTable("net"+strconv.Itoa(Nex)+".txt", colRowTable)
	Nex++
	algX4(&colRowTable)
	return board
}

func algX4(colRowTable *[64]uint64) {
	//проверим isTableEmpty

	numOneInCol, numOneRow, success := findOneInCol(colRowTable)
	if success {
		fmt.Println("numOneInCol = ", numOneInCol, ", numOneRow = ", numOneRow)
		coverRowsWithOnes(numOneInCol, numOneRow, colRowTable)
		printTable("net"+strconv.Itoa(Nex)+".txt", *colRowTable)
		Nex++
	} else {
		fmt.Println("can't find one in col")
	}

}

func coverRowsWithOnes(numOneInCol, numOneRow int, colRowTable *[64]uint64) { //накрываем соответствующие строки
	//в строке с одной единицей найдем столбцы где есть единица
	var onesSlice []int = make([]int, 0) //слайс с номерами столбцов где единица
	for i := 0; i < 64; i++ {            //проходим по массиву чисел
		mulCheck := colRowTable[numOneRow]
		mulCheck |= 1 << i                      //В i-том разряде установим единицу
		if mulCheck == colRowTable[numOneRow] { //если равны значит в этом разряде единица
			onesSlice = append(onesSlice, i)
		}
	}
	fmt.Println("onesSlice = ", onesSlice)
	//для опорной строки нашли столбцы где единица, в этих столбцах будем находить строки где в этом разряде единица
	var numColOne []uint64 = make([]uint64, 0)
	var rowOnesSlice []int = make([]int, 0) //слайс с номерами строк где нашли единицу
	for _, col := range onesSlice {         //создадим слайс из чисел с единицей в соответствующем разряде
		var mul uint64
		mul |= 1 << col
		numColOne = append(numColOne, mul)
	}
	//fmt.Println(numColOne)
	/*for _, mul := range numColOne {
		fmt.Printf("%064b\n", mul)
	}*/
	for i := 0; i < 64; i++ {
		for _, mul := range numColOne {
			mulCheck := colRowTable[i]
			mulCheck |= mul
			if mulCheck == colRowTable[i] {
				rowOnesSlice = append(rowOnesSlice, i)
				break
			}
		}
	}
	fmt.Println("rowOnesSlice = ", rowOnesSlice)
	for _, numRowCower := range rowOnesSlice {
		colRowTable[numRowCower] = 0
	}
}

func findOneInCol(colRowTable *[64]uint64) (int, int, bool) { //найдем столбец с одной единицей и строку где эта единица
	var cntOnes, rowOne int
	//Получим копию i-го числа массива colRowTable. В j-том разряде установим единицу
	//если число не изменится, то в этом разряде у него единица
	for j := 0; j < 64; j++ { //проходим по столбцу (номер разряда)
		cntOnes = 0
		for i := 0; i < 64; i++ { //проходим по массиву чисел
			if colRowTable[i] == 0 {
				continue
			}
			mulCheck := colRowTable[i]
			mulCheck |= 1 << j //В j-том разряде установим единицу
			if mulCheck == colRowTable[i] {
				//fmt.Printf("%064b\n", colRowTable[i])
				//fmt.Printf("%064b\n", mulCheck)
				rowOne = i
				cntOnes++
			}
		}
		//fmt.Println("cntOnes = ", cntOnes)
		if cntOnes == 1 {
			return j, rowOne, true
		}
	}
	return 0, 0, false
}

func fillTable(board *[4][4]int, colRowTable *[64]uint64) /*[64]uint64*/ {
	//var colRowTable [64]uint64
	var mul1 uint64 = 1
	mul2 := [4]uint64{1 << 16, 1 << 20, 1 << 24, 1 << 28}
	mul3 := [4]uint64{1 << 32, 1 << 36, 1 << 40, 1 << 44}
	mul4 := [4]uint64{1 << 48, 1 << 52, 1 << 56, 1 << 60}

	//mul1 = 1
	for j := 0; j < 16; j++ {
		for i := 0; i < 4; i++ {
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul1    //ограничения в ячейках
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul2[i] //ограничения в строках
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul3[i] //ограничения в столбцах
			colRowTable[j*4+i] = colRowTable[j*4+i] | mul4[i] //ограничения в боксах
		}
		mul1 = mul1 << 1

		if j == 3 || j == 7 || j == 11 {
			for k := 0; k < 4; k++ {
				mul2[k] = mul2[k] << 1
			}
		}

		for k := 0; k < 4; k++ {
			mul3[k] = mul3[k] << 1
		}
		if j == 3 || j == 7 || j == 11 {
			for k := 0; k < 4; k++ {
				mul3[k] = mul3[k] >> 4
			}
		}

		if j == 1 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] << 1
			}
		}
		if j == 3 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] >> 1
			}
		}
		if j == 5 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] << 1
			}
		}
		if j == 7 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] << 1
			}
		}
		if j == 9 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] << 1
			}
		}
		if j == 11 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] >> 1
			}
		}
		if j == 13 {
			for k := 0; k < 4; k++ {
				mul4[k] = mul4[k] << 1
			}
		}
	}

	//сформировали colRowTable как со всеми незаполненными ячейками board
	//накроем строки в colRowTable для заполненных ячеек board
	var row, col int
	var mul5 uint64
	var endLoop bool
	for i := 0; i < 64; i += 4 {
		if board[row][col] == 0 {
			row, col, endLoop = nextCellBoard(row, col)
			if endLoop {
				break
			} else {
				continue
			}
		}
		//накрываем строки кроме значения в ячейке board
		for k := 1; k <= 4; k++ {
			if board[row][col] == k {
				continue
			} else {
				colRowTable[i+k-1] = colRowTable[i+k-1] & mul5
			}
		}

		row, col, endLoop = nextCellBoard(row, col)
		if endLoop {
			break
		}
	}
}

func nextCellBoard(row, col int) (int, int, bool) {
	var endLoop bool = false
	if col == 3 {
		col = 0
		row++
	} else {
		col++
	}
	if row == 3 && col == 3 {
		endLoop = true
	}
	return row, col, endLoop
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

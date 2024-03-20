package main

import (
	"fmt"
	"os"
	"sync"
)

type Node struct { //обычный узел
	data      int
	row       int
	col       int
	nextRight *Node
	nextLeft  *Node
	nextUp    *Node
	nextDown  *Node
}

type List struct { //список с узлами
	head *Node
}

func main() {
	board := [4][4]int{} //board[row][col]
	board[0][3] = 3
	board[1][0] = 4
	board[2][1] = 1
	board[2][2] = 3
	board[3][0] = 3
	board[3][2] = 2
	fmt.Println(board)
	SolveAlgX(board)
}

func SolveAlgX(board [4][4]int) [4][4]int {
	var l = List{}
	l = fillHeads(l)
	printListToFile(l, "output.txt")
	t := createT()
	printListToFile(t, "output2.txt")
	doRestrict(board, l, t)
	printListToFile(l, "output3.txt")
	return [4][4]int{}
}

func doRestrict(board [4][4]int, l, t List) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if board[row][col] != 0 {
				//numRowTableRestr := 4*row + 2*col + board[row][col] //нашли какая строка должна остаться
				fmt.Println(row, col)
				numsRowRemove := findRemove(row, col, board)
				fmt.Println(numsRowRemove)
				l = doCowerRows(l, numsRowRemove) //сделаем накрытие строки numRowRemove
			}
		}
	}
	solveSudoku(l, t)
}

func solveSudoku(l, t List) {

	listInCols1 := findInCols1(l)
	fmt.Println(listInCols1)
}

func findInCols1(l List) []int {
	var listInCols1 []int = make([]int, 0)
	colMain := l.head.nextRight
	for i := 1; i <= 64; i++ {
		if colMain.nextDown != nil && colMain.nextDown.nextDown == nil {
			listInCols1 = append(listInCols1, i)
		}
		if i < 64 {
			colMain = colMain.nextRight
		}
	}

	return listInCols1
}

func findRemove(row, col int, board [4][4]int) []int {
	var numsRowRemove []int = make([]int, 0)
	numRowTableRestr := 16*row + 4*col + board[row][col] //нашли какая строка должна остаться
	fmt.Println(numRowTableRestr)
	jumpToStart := 16*row + 4*col
	for i := 1; i <= 4; i++ {
		if i+jumpToStart != numRowTableRestr {
			numsRowRemove = append(numsRowRemove, i+jumpToStart)
		}
	}
	return numsRowRemove
}

func doCowerRows(l List, numsRowRemove []int) List {
	for _, numRow := range numsRowRemove {
		rowCower := l.head
		for rowCower.row != numRow {
			rowCower = rowCower.nextDown
		}
		for {
			if rowCower.nextDown == nil {
				//rowCower.nextUp.nextDown = nil
				upCell := rowCower.nextUp
				upCell.nextDown = nil
			} else {
				//rowCower.nextUp.nextDown = rowCower.nextDown
				//rowCower.nextDown.nextUp = rowCower.nextUp
				upCell := rowCower.nextUp
				downCell := rowCower.nextDown
				upCell.nextDown = downCell
				downCell.nextUp = upCell
			}

			if rowCower.nextRight != nil {
				rowCower = rowCower.nextRight
			} else {
				break
			}
		}
	}
	return l
}

func createT() List {
	t := List{}
	headList := &Node{data: 1, col: 0, row: 0}
	t.head = headList
	currNode := t.head
	for i := 1; i <= 64; i++ {
		newNodeCol := &Node{data: 1, col: i, row: 0}
		currNode.nextRight = newNodeCol
		newNodeCol.nextLeft = currNode
		currNode = newNodeCol
	}
	return t
}

func fillHeads(l List) List {
	headList := &Node{data: 1, col: 0, row: 0}
	l.head = headList
	var wg sync.WaitGroup
	wg.Add(2)
	//Создаем заголовочные ноды столбцов
	go func() {
		currNode := l.head
		for i := 1; i <= 64; i++ {
			newNodeCol := &Node{data: 1, col: i, row: 0}
			currNode.nextRight = newNodeCol
			newNodeCol.nextLeft = currNode
			currNode = newNodeCol
		}
		currNode.nextRight = l.head //соединяем вкруговую
		l.head.nextLeft = currNode
		wg.Done()
	}()
	//создаем заголовочные ноды строк
	go func() {
		currNode := l.head
		for i := 1; i <= 64; i++ {
			newNodeCol := &Node{data: 1, col: 0, row: i}
			currNode.nextDown = newNodeCol
			newNodeCol.nextUp = currNode
			currNode = newNodeCol
		}
		currNode.nextDown = l.head //соединяем вкруговую
		l.head.nextUp = currNode
		wg.Done()
	}()
	wg.Wait()
	//--- Заполняем ограничения в ячейках
	wg.Add(16)
	for rc := 0; rc < 16; rc++ {
		go func(row, col int) {
			currNodeDown := l.head
			currNodeRight := l.head

			for currNodeDown.row != row { //опускаемся на нужную строку
				currNodeDown = currNodeDown.nextDown
			}
			for i := 0; i < col; i++ { //передвигаемся вправо на нужный столбец
				currNodeRight = currNodeRight.nextRight
			}

			for r := row; r < row+4; r++ {
				newNode := &Node{data: 1, col: col, row: r}
				currNodeDown.nextRight = newNode //привязываемся к новой ноде справа
				newNode.nextLeft = currNodeDown

				newNode.nextUp = currNodeRight
				currNodeRight.nextDown = newNode

				currNodeRight = newNode
				currNodeDown = currNodeDown.nextDown
			}
			wg.Done()
		}(rc*4+1, rc+1)
	}
	wg.Wait()
	//--- Заполняем ограничения в строках
	rc := 17
	rr := 1
	wg.Add(16)
	for irs := 0; irs < 4; irs++ {
		for ir := 0; ir < 4; ir++ {
			go func(row, col int) {
				currNodeDown := l.head
				currNodeRight := l.head

				for currNodeDown.row != row { //опускаемся на нужную строку
					currNodeDown = currNodeDown.nextDown
				}
				for i := 0; i < col; i++ { //передвигаемся вправо на нужный столбец
					currNodeRight = currNodeRight.nextRight
				}

				for i := 0; i < 4; i++ {
					newNode := &Node{data: 1, col: col, row: row + i*4}
					currNodeDown.nextRight.nextRight = newNode //привязываемся к новой ноде справа
					newNode.nextLeft = currNodeDown.nextRight

					newNode.nextUp = currNodeRight
					currNodeRight.nextDown = newNode

					currNodeRight = newNode
					for j := 0; j < 4; j++ {
						currNodeDown = currNodeDown.nextDown
					}
				}
				wg.Done()
			}(rr+16*ir+irs, rc+ir+irs*4)
		}
	}
	wg.Wait()
	//--- Заполняем ограничения в столбцах
	rc = 33
	rr = 1
	wg.Add(16)
	for irs := 0; irs < 4; irs++ {
		for ir := 0; ir < 4; ir++ {
			go func(row, col int) {
				currNodeDown := l.head
				currNodeRight := l.head

				for currNodeDown.row != row { //опускаемся на нужную строку
					currNodeDown = currNodeDown.nextDown
				}
				for i := 0; i < col; i++ { //передвигаемся вправо на нужный столбец
					currNodeRight = currNodeRight.nextRight
				}

				for i := 0; i < 4; i++ {
					newNode := &Node{data: 1, col: col, row: row + i*16}
					currNodeDown.nextRight.nextRight.nextRight = newNode //привязываемся к новой ноде справа
					newNode.nextLeft = currNodeDown.nextRight.nextRight

					newNode.nextUp = currNodeRight
					currNodeRight.nextDown = newNode

					currNodeRight = newNode
					for j := 0; j < 16; j++ {
						currNodeDown = currNodeDown.nextDown
					}
				}
				wg.Done()
			}(rr+4*ir+irs, rc+ir+irs*4)
		}
	}
	wg.Wait()
	//--- Заполняем ограничения в боксах
	rc = 49
	rr = 1
	wg.Add(16)
	for irs := 0; irs < 4; irs++ {
		irb := 0
		switch irs {
		case 1:
			irb = 8
		case 2:
			irb = 32
		case 3:
			irb = 40
		}

		for ir := 0; ir < 4; ir++ {
			go func(row, col int) {
				currNodeDown := l.head
				currNodeRight := l.head

				for currNodeDown.row != row { //опускаемся на нужную строку
					currNodeDown = currNodeDown.nextDown
				}
				for i := 0; i < col; i++ { //передвигаемся вправо на нужный столбец
					currNodeRight = currNodeRight.nextRight
				}

				for i := 0; i < 4; i++ {
					switch i {
					case 0:
						{
							newNode := &Node{data: 1, col: col, row: row}
							currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode //привязываемся к новой ноде справа
							newNode.nextLeft = currNodeDown.nextRight.nextRight.nextRight
							newNode.nextUp = currNodeRight
							currNodeRight.nextDown = newNode
							for j := 0; j < 4; j++ {
								currNodeDown = currNodeDown.nextDown
							}
						}
					case 1:
						{
							newNode := &Node{data: 1, col: col, row: row + 4}
							currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode //привязываемся к новой ноде справа
							newNode.nextLeft = currNodeDown.nextRight.nextRight.nextRight
							newNode.nextUp = currNodeRight
							currNodeRight.nextDown = newNode
							for j := 0; j < 12; j++ {
								currNodeDown = currNodeDown.nextDown
							}
						}
					case 2:
						{
							newNode := &Node{data: 1, col: col, row: row + 16}
							currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode //привязываемся к новой ноде справа
							newNode.nextLeft = currNodeDown.nextRight.nextRight.nextRight
							newNode.nextUp = currNodeRight
							currNodeRight.nextDown = newNode
							for j := 0; j < 4; j++ {
								currNodeDown = currNodeDown.nextDown
							}
						}
					case 3:
						{
							newNode := &Node{data: 1, col: col, row: row + 20}
							currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode //привязываемся к новой ноде справа
							newNode.nextLeft = currNodeDown.nextRight.nextRight.nextRight
							newNode.nextUp = currNodeRight
							/*currNodeRight.nextDown = newNode
							for j := 0; j < 4; j++ {
								currNodeDown = currNodeDown.nextDown
							}*/
						}
					}

				}
				wg.Done()
			}(rr+ir+irb, rc+4*ir+irs)
		}
	}
	wg.Wait()
	return l
}

func printListToFile(l List, fileName string) {
	currNode := l.head
	// open output file
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer fo.Close()

	for i := 0; i <= 64; i++ {
		//fmt.Printf("%d\t", currNode.col)
		_, err = fo.WriteString(fmt.Sprintf("%d\t", i)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		currNode = currNode.nextRight
	}
	_, err = fo.WriteString(fmt.Sprintf("\n")) // writing...
	if err != nil {
		fmt.Printf("error writing string: %v", err)
	}
	currNode = l.head.nextDown

	for currNode != l.head {
		if currNode == nil {
			break
		}
		//fmt.Printf("%d", currNode.row) //заголовок строки
		_, err = fo.WriteString(fmt.Sprintf("%d", currNode.row)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		amountTabs := 0 //номер предыдущего столбца
		currNodeRow := currNode
		for currNodeRow.nextRight != nil {
			amountTabs = currNodeRow.nextRight.col - currNodeRow.col
			for k := 0; k < amountTabs; k++ {
				//fmt.Printf("\t")
				_, err = fo.WriteString(fmt.Sprintf("\t")) // writing...
				if err != nil {
					fmt.Printf("error writing string: %v", err)
				}
			}
			//fmt.Printf("1")
			_, err = fo.WriteString(fmt.Sprintf("1")) // writing...
			if err != nil {
				fmt.Printf("error writing string: %v", err)
			}
			currNodeRow = currNodeRow.nextRight
		}
		//fmt.Printf("\n")
		_, err = fo.WriteString(fmt.Sprintf("\n")) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		currNode = currNode.nextDown
	}
}

func printList(l List) {
	currNode := l.head
	for i := 0; i <= 64; i++ {
		fmt.Printf("%d\t", currNode.col)
		currNode = currNode.nextRight
	}
	fmt.Printf("\n") //распечатали заголовки столбцов
}

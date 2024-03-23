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
	//listRows := checkCol(l, 64)
	//fmt.Println(listRows)
	doRestrict(board, l, t)
	printListToFile(l, "output3.txt")
	return [4][4]int{}
}

func doRestrict(board [4][4]int, l, t List) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if board[row][col] != 0 {
				fmt.Println(row, col)
				numsRowRemove := findRemove(row, col, board)
				fmt.Println(numsRowRemove)
				l = doCowerRows(l, numsRowRemove) //сделаем накрытие строки numRowRemove
			}
		}
	}
	listInCols1 := findInCols1(l)
	listInCols1 = listInCols1[0:]
	fmt.Println(listInCols1)
	solveSudoku(l, t, listInCols1)
}

func solveSudoku(l, t List, listInCols1 []int) {
	var numRowWithOne int //номер строки с одной единицей в столбце
	var arrColsOne []int
	if len(listInCols1) > 0 {
		numRowWithOne, arrColsOne = findNumRowWithOne(l, listInCols1[0])
		fmt.Println(numRowWithOne)
	}
	t = addToTableProbableSolution(numRowWithOne, l, t)
	printListToFile(t, "output4.txt")
	fmt.Println(arrColsOne)
	//listRows := checkCol(l, 61)
	//fmt.Println(listRows)
}

func addToTableProbableSolution(numRowWithOne int, l, t List) List { //будем добавлять строку в таблицу возможных ответов
	rowMainL := l.head
	for rowMainL.row != numRowWithOne {
		rowMainL = rowMainL.nextDown
	}
	rowMainT := t.head
	for rowMainT.nextDown != nil {
		rowMainT = rowMainT.nextDown
	}
	newNodeCol := &Node{data: 1, col: 0, row: numRowWithOne}
	rowMainT.nextDown = newNodeCol
	newNodeCol.nextUp = rowMainT
	for rowMainL.nextRight != nil {
		rowMainL = rowMainL.nextRight
		addNode(t, rowMainL.col, rowMainL.row)
	}
	return t
}

func addNode(l List, col, row int) List {
	currNodeDown := l.head
	currNodeRight := l.head
	newNode := &Node{data: 1, col: col, row: row}
	for currNodeDown.row != row { //опускаемся на нужную строку
		currNodeDown = currNodeDown.nextDown
	}
	for currNodeDown.nextRight != nil { //в этой строке передвигаемся в крайнее правое положение
		currNodeDown = currNodeDown.nextRight
	}
	currNodeDown.nextRight = newNode //привязываемся к новой ноде справа
	newNode.nextLeft = currNodeDown

	for i := 0; i < col; i++ { //передвигаемся вправо на нужный столбец
		currNodeRight = currNodeRight.nextRight
	}
	for currNodeRight.nextDown != nil { //передвигаемся в столбце в крайнее нижнее положение
		currNodeRight = currNodeRight.nextDown
	}

	newNode.nextUp = currNodeRight
	currNodeRight.nextDown = newNode
	return l
}

func findNumRowWithOne(l List, colNum int) (int, []int) { //находим номер строки где в столбце одна единица и список столбцов с единицами для закрытия
	var numRowOne int
	var arrColsOne []int = make([]int, 0)
	colMain := l.head
	for colMain.col != colNum {
		colMain = colMain.nextRight
	}
	colMain = colMain.nextDown
	numRowOne = colMain.row
	//составим список столбцов с единицами
	for colMain.nextLeft != nil {
		colMain = colMain.nextLeft
	}
	for colMain.nextRight != nil {
		colMain = colMain.nextRight
		arrColsOne = append(arrColsOne, colMain.col)
	}
	return numRowOne, arrColsOne
}

func checkCol(l List, col int) []int {
	var listInCols1 []int = make([]int, 0)
	colMain := l.head
	for i := 0; i < 64; i++ {
		if colMain.nextRight != nil {
			colMain = colMain.nextRight
		}
		if colMain.col == col {
			fmt.Println(colMain.col)
			for {
				listInCols1 = append(listInCols1, colMain.row)
				if colMain.nextDown != nil {
					colMain = colMain.nextDown
				} else {
					break
				}

			}

			break
		}
	}
	return listInCols1
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
	//--- Заполняем ограничения в боксах Доделать!!!!!!!!
	rc = 49
	rr = 1

	for inb := 0; inb < 4; inb++ {

		irbc := []int{0, 8, 32, 40}
		wg.Add(4)
		for ir := 0; ir < 4; ir++ {
			go func(row, col int) {
				currNodeDown := l.head
				currNodeRight := l.head
				for i := 0; i < row; i++ {
					currNodeDown = currNodeDown.nextDown
				}
				for i := 0; i < col; i++ {
					currNodeRight = currNodeRight.nextRight
				} //поставили заголовочные ноды на места

				irb := []int{0, 4, 16, 20}
				for irs := 0; irs < 4; irs++ {
					newNode := &Node{data: 1, col: col, row: row + irb[irs]}
					currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode

					currNodeRight.nextDown = newNode
					newNode.nextUp = currNodeRight
					currNodeRight = newNode

					if irs < 3 {
						for i := 0; i < irb[irs+1]-irb[irs]; i++ {
							currNodeDown = currNodeDown.nextDown
						}
					}
				}
				wg.Done()
			}(rr+irbc[ir], rc+ir)
		}
		wg.Wait()
		rc += 4
		rr += 1
	}
	/*go func(row, col int) {
		currNodeDown := l.head
		currNodeRight1 := l.head
		currNodeRight2 := l.head
		currNodeRight3 := l.head
		currNodeRight4 := l.head
		wg.Add(5)
		go func() {
			for i := 0; i < row; i++ {
				currNodeDown = currNodeDown.nextDown
			}
			wg.Done()
		}()
		go func() {
			for i := 0; i < col; i++ {
				currNodeRight1 = currNodeRight1.nextRight
			}
			for currNodeRight1.nextDown != nil {
				currNodeRight1 = currNodeRight1.nextDown
			}
			wg.Done()
		}()
		go func() {
			for i := 0; i < col+1; i++ {
				currNodeRight2 = currNodeRight2.nextRight
			}
			for currNodeRight2.nextDown != nil {
				currNodeRight2 = currNodeRight2.nextDown
			}
			wg.Done()
		}()
		go func() {
			for i := 0; i < col+2; i++ {
				currNodeRight3 = currNodeRight3.nextRight
			}
			for currNodeRight3.nextDown != nil {
				currNodeRight3 = currNodeRight3.nextDown
			}
			wg.Done()
		}()
		go func() {
			for i := 0; i < col+3; i++ {
				currNodeRight4 = currNodeRight4.nextRight
			}
			for currNodeRight4.nextDown != nil {
				currNodeRight4 = currNodeRight4.nextDown
			}
			wg.Done()
		}()
		wg.Wait() //поставили заголовочные ноды на места

		irb := []int{0, 8, 32, 40}
		//for irt := 0; irt < 2; irt++ {
		for irk := 0; irk < 2; irk++ {
			for irs := 0; irs < 4; irs++ {
				newNode := &Node{data: 1, col: col + irs, row: row + irb[irs]}
				currNodeDown.nextRight.nextRight.nextRight.nextRight = newNode
				switch irs {
				case 0:
					currNodeRight1.nextDown = newNode
					newNode.nextUp = currNodeRight1
					currNodeRight1 = newNode
				case 1:
					currNodeRight2.nextDown = newNode
					newNode.nextUp = currNodeRight2
					currNodeRight2 = newNode
				case 2:
					currNodeRight3.nextDown = newNode
					newNode.nextUp = currNodeRight3
					currNodeRight3 = newNode
				case 3:
					currNodeRight4.nextDown = newNode
					newNode.nextUp = currNodeRight4
					currNodeRight4 = newNode
				}
				if irs < 3 {
					for i := 0; i < irb[irs+1]-irb[irs]; i++ {
						currNodeDown = currNodeDown.nextDown
					}
				}
			}
			for i := 0; i < 36; i++ {
				currNodeDown = currNodeDown.nextUp
			}
			row += 4
		}
		//	for i := 0; i < 28; i++ {
		//		currNodeDown = currNodeDown.nextUp
		//	}
		//	row += 12
		//}
	}(rr, rc)*/

	//wg.Add(16)
	/*for irs := 0; irs < 4; irs++ {
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
		fillBoxes(rr+ir+irb, rc+4*ir+irs, l)*/
	/*go func(row, col int) {
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
				}
			}

		}
		wg.Done()
	}(rr+ir+irb, rc+4*ir+irs)*/
	//}
	//}
	//wg.Wait()
	return l
}

func fillBoxes(row, col int, l List) {
	//go func(row, col int) {
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
			}
		}

	}
	//wg.Done()
	//}(rr+ir+irb, rc+4*ir+irs)
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

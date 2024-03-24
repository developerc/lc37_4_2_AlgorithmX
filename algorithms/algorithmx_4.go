package algorithms

import (
	"fmt"
	"os"
	"sort"
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

//var loops int

func SolveAlgX(board [4][4]int) [4][4]int {
	var l = List{}
	l = fillHeads(l)
	//printListToFile(l, "output.txt")
	t := createT()
	//printListToFile(t, "output2.txt")
	doRestrict(board, l, t)
	//здесь из t извлекаем ответы
	result := findResultBoard(t)
	return result
}

func findResultBoard(t List) [4][4]int {
	board := [4][4]int{} //board[row][col]
	var val int
	var r int
	var c int
	rowHead := t.head
	for rowHead.nextDown != nil {
		rowHead = rowHead.nextDown
		//определим значение
		if rowHead.row%2 == 0 && rowHead.row%4 == 0 {
			val = 4
		}
		if rowHead.row%2 == 0 && rowHead.row%4 != 0 {
			val = 2
		}
		if rowHead.row == 1 {
			val = 1
		}
		if rowHead.row == 3 {
			val = 3
		}
		if rowHead.row != 1 && rowHead.row != 3 && (rowHead.row-1)%2 == 0 && (rowHead.row-1)%4 != 0 {
			val = 3
		}
		if rowHead.row != 1 && rowHead.row != 3 && (rowHead.row-1)%2 == 0 && (rowHead.row-1)%4 == 0 {
			val = 1
		}
		// определим row
		r = rowHead.row / 16
		if rowHead.row%16 == 0 {
			r--
		}
		/*if rowHead.row >= 1 && rowHead.row <= 16 {
			r = 0
		}
		if rowHead.row >= 17 && rowHead.row <= 32 {
			r = 1
		}
		if rowHead.row >= 33 && rowHead.row <= 48 {
			r = 2
		}
		if rowHead.row >= 49 && rowHead.row <= 64 {
			r = 3
		}*/
		//определим col
		if rowHead.row%16 != 0 {
			var cd int = rowHead.row % 16
			c = cd / 4
			if cd%4 == 0 {
				c--
			}
			//fmt.Printf("rowHead.row = %d, cd = %d, r = %d, c = %d, val = %d\n", rowHead.row, cd, r, c, val)
		} else {
			c = 3
		}
		board[r][c] = val
	}

	return board
}

func doRestrict(board [4][4]int, l, t List) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if board[row][col] != 0 {
				//fmt.Println(row, col)
				numsRowRemove := findRemove(row, col, board)
				//fmt.Println(numsRowRemove)
				l = doCowerRows(l, numsRowRemove) //сделаем накрытие строки numRowRemove
			}
		}
	}
	solveSudoku(l, t)
}

func solveSudoku(l, t List) {
	/*if loops > 1 {
		return
	}*/
	if listIsEmpty(l) {
		return
	}
	var numRowWithOne int         //номер строки с одной единицей в столбце
	var arrColsOne []int          //список столбцов в строке где в столбце одна единица
	listInCols1 := findInCols1(l) //составляем список столбцов с одной единицей
	listInCols1 = listInCols1[0:]
	if len(listInCols1) > 0 {
		numRowWithOne, arrColsOne = findNumRowWithOne(l, listInCols1[0])
		//fmt.Println(numRowWithOne)
	}
	t = addToTableProbableSolution(numRowWithOne, l, t) //в таблицу возможных решений добавляем строку
	//printListToFile(t, "output4.txt")
	//fmt.Println(arrColsOne)
	arrRowsToCover := findArrRowsToCover(arrColsOne, l)
	//fmt.Println(arrRowsToCover)
	coverRows(arrRowsToCover, l)
	//printListToFile(l, "output5.txt")
	//loops++
	solveSudoku(l, t)
}

func listIsEmpty(l List) bool {
	return l.head.nextDown == l.head
}

func coverRows(arrRowsToCover []int, l List) {
	if len(arrRowsToCover) == 0 {
		return
	}
	m, uniq := make(map[int]struct{}), make([]int, 0, len(arrRowsToCover)) //убираем повторяющиеся номера
	for _, v := range arrRowsToCover {
		if _, ok := m[v]; !ok {
			m[v], uniq = struct{}{}, append(uniq, v)
		}
	}
	sort.Ints(uniq)
	//fmt.Println(uniq)
	// делаем накрытие строк
	rowCower := l.head
	for i := 0; i < len(uniq); i++ {
		for rowCower.row != uniq[i] {
			rowCower = rowCower.nextDown
		}
		rowTemp := rowCower
		for {
			if rowTemp.nextDown == nil {
				rowTemp.nextUp.nextDown = nil
			} else {
				rowTemp.nextUp.nextDown = rowTemp.nextDown
				rowTemp.nextDown.nextUp = rowTemp.nextUp
			}

			if rowTemp.nextRight != nil {
				rowTemp = rowTemp.nextRight
			} else {
				break
			}
		}
	}
}

func findArrRowsToCover(arrColsOne []int, l List) []int {
	var arrRowsToCover []int = make([]int, 0)
	colMain := l.head
	for _, col := range arrColsOne {
		for colMain.col != col {
			colMain = colMain.nextRight
		}
		//fmt.Println("col for cover: ", col)
		nodeForCover := colMain.nextDown
		arrRowsToCover = append(arrRowsToCover, nodeForCover.row)
		for nodeForCover.nextDown != nil {
			nodeForCover = nodeForCover.nextDown
			arrRowsToCover = append(arrRowsToCover, nodeForCover.row)
		}
	}
	return arrRowsToCover
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
	//fmt.Println(numRowTableRestr)
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
				rowCower.nextUp.nextDown = nil
			} else {
				rowCower.nextUp.nextDown = rowCower.nextDown
				rowCower.nextDown.nextUp = rowCower.nextUp
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

	return l
}

func printListToFile(l List, fileName string) {
	currNode := l.head
	fo, err := os.Create(fileName) // open output file
	if err != nil {
		panic(err)
	}
	defer fo.Close() // close fo on exit and check for its returned error

	for i := 0; i <= 64; i++ {
		_, err = fo.WriteString(fmt.Sprintf("%d\t", i)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		currNode = currNode.nextRight
	}
	_, err = fo.WriteString("\n") // writing...
	if err != nil {
		fmt.Printf("error writing string: %v", err)
	}
	currNode = l.head.nextDown

	for currNode != l.head {
		if currNode == nil {
			break
		}
		_, err = fo.WriteString(fmt.Sprintf("%d", currNode.row)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		amountTabs := 0 //номер предыдущего столбца
		currNodeRow := currNode
		for currNodeRow.nextRight != nil {
			amountTabs = currNodeRow.nextRight.col - currNodeRow.col
			for k := 0; k < amountTabs; k++ {
				_, err = fo.WriteString("\t") // writing...
				if err != nil {
					fmt.Printf("error writing string: %v", err)
				}
			}
			_, err = fo.WriteString("1") // writing...
			if err != nil {
				fmt.Printf("error writing string: %v", err)
			}
			currNodeRow = currNodeRow.nextRight
		}
		_, err = fo.WriteString("\n") // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
		currNode = currNode.nextDown
	}
}

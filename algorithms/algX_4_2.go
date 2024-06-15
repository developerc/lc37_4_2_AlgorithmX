package algorithms

import (
	"fmt"
	"os"
	"strconv"
)

type NodeNet struct {
	colNet [64]Node
	rowNet [64]Node
}

var n int

func SolveAlgX4_2(board [4][4]int) [4][4]int {
	colRowNet := fillNet(&board)
	printNetToFile("net"+strconv.Itoa(n)+".txt", colRowNet)
	algX(colRowNet, &board)
	return board
}

func algX(colRowNet NodeNet, board *[4][4]int) [4][4]int {

	return *board
}

func fillNet(board *[4][4]int) NodeNet {
	nodeNet := NodeNet{}
	for i := 0; i < 64; i++ {
		nodeNet.rowNet[i].col = i
	}

	var cNet, rNet int
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			valEnd := 1
			if board[row][col] == 0 { //если val равен 0 то добавляем Node для 1, 2, 3, 4
				valEnd = 4
			}

			//ограничения в ячейках
			cNet = col + 4*row //столбец
			for val := 1; val <= valEnd; val++ {
				if valEnd == 1 {
					rNet = 4*col + 16*row + board[row][col] - 1
				} else {
					rNet = 4*col + 16*row + val - 1 //ряд
				}
				node := &Node{row: rNet, col: cNet, data: 1}
				nodeNet.colNet[rNet].nextRight = node
				node.nextLeft = &nodeNet.colNet[rNet]
				nodeTemp := nodeNet.rowNet[cNet]
				for nodeTemp.nextDown != nil {
					nodeTemp = *nodeTemp.nextDown
				}
				nodeTemp.nextDown = node
				node.nextUp = &nodeTemp
			}
			//ограничения в строках
			/*for val := 1; val <= valEnd; val++ {
				cNet = row + (val-1)*4 + 16 //столбец
				for val := 1; val <= valEnd; val++ {
					if valEnd == 1 {
						rNet = 4*col + 16*row + board[row][col] - 1
					} else {
						rNet = 4*col + 16*row + val - 1 //ряд
					}
					node := &Node{row: rNet, col: cNet, data: 1}
					nodeNet.colNet[rNet].nextRight.nextRight = node
					node.nextLeft = nodeNet.colNet[rNet].nextRight
					nodeTemp := nodeNet.rowNet[cNet]
					for nodeTemp.nextDown != nil {
						nodeTemp = *nodeTemp.nextDown
					}
					nodeTemp.nextDown = node
					node.nextUp = &nodeTemp
				}

			}*/

		}
	}

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			valEnd := 1
			if board[row][col] == 0 { //если val равен 0 то добавляем Node для 1, 2, 3, 4
				valEnd = 4
			}
			//ограничения в строках
			for val := 1; val <= valEnd; val++ {
				cNet = row + (val-1)*4 + 16 //столбец
				for val := 1; val <= valEnd; val++ {
					if valEnd == 1 {
						rNet = 4*col + 16*row + board[row][col] - 1
					} else {
						rNet = 4*col + 16*row + val - 1 //ряд
					}
					node := &Node{row: rNet, col: cNet, data: 1}
					nodeNet.colNet[rNet].nextRight.nextRight = node
					node.nextLeft = nodeNet.colNet[rNet].nextRight
					nodeTemp := nodeNet.rowNet[cNet]
					for nodeTemp.nextDown != nil {
						nodeTemp = *nodeTemp.nextDown
					}
					nodeTemp.nextDown = node
					node.nextUp = &nodeTemp
				}

			}
		}
	}
	return nodeNet
}

func printNetToFile(fileName string, colRowNet NodeNet) {
	fo, err := os.Create(fileName) // open output file
	if err != nil {
		panic(err)
	}
	defer fo.Close() // close fo on exit and check for its returned error
	for i := 0; i < 64; i++ {
		_, err = fo.WriteString(fmt.Sprintf("%d\t", i)) // пишем заголовок столбцов
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
	for i := 0; i < 64; i++ {
		fo.WriteString("\n")
		if colRowNet.colNet[i].nextRight == nil {
			continue
		}
		numTabs := colRowNet.colNet[i].nextRight.col
		for k := 0; k < numTabs; k++ {
			fo.WriteString("\t")
		}
		fo.WriteString("1")
		//--
		numTabs = colRowNet.colNet[i].nextRight.nextRight.col - colRowNet.colNet[i].nextRight.col
		for k := 0; k < numTabs; k++ {
			fo.WriteString("\t")
		}
		fo.WriteString("1")
		//--
	}

}

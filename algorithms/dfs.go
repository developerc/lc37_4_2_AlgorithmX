package algorithms

func SolveDfs(board [4][4]int) [4][4]int {
	backtrack(&board)
	return board
}

func backtrack(board *[4][4]int) bool {
	if !hasEmptyCell(board) {
		return true
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				for candidate := 4; candidate >= 1; candidate-- {
					board[i][j] = candidate
					if isBoardValid(board) {
						if backtrack(board) {
							return true
						}
						board[i][j] = 0
					} else {
						board[i][j] = 0
					}
				}
				return false
			}
		}
	}
	return false
}

func hasEmptyCell(board *[4][4]int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				return true
			}
		}
	}
	return false
}

func isBoardValid(board *[4][4]int) bool {

	//check duplicates by row
	for row := 0; row < 4; row++ {
		counter := [5]int{}
		for col := 0; col < 4; col++ {
			counter[board[row][col]]++
		}
		if hasDuplicates(counter) {
			return false
		}
	}

	//check duplicates by column
	for row := 0; row < 4; row++ {
		counter := [5]int{}
		for col := 0; col < 4; col++ {
			counter[board[col][row]]++
		}
		if hasDuplicates(counter) {
			return false
		}
	}

	//check 2x2 section
	for i := 0; i < 4; i += 2 {
		for j := 0; j < 4; j += 2 {
			counter := [5]int{}
			for row := i; row < i+2; row++ {
				for col := j; col < j+2; col++ {
					counter[board[row][col]]++
				}
				if hasDuplicates(counter) {
					return false
				}
			}
		}
	}

	return true
}

func hasDuplicates(counter [5]int) bool {
	for i, count := range counter {
		if i == 0 {
			continue
		}
		if count > 1 {
			return true
		}
	}
	return false
}

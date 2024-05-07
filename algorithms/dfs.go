package algorithms

func backtrack(board [4][4]int) bool {
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

func hasEmptyCell(board [4][4]int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if board[i][j] == 0 {
				return true
			}
		}
	}
	return false
}

func isBoardValid(board [4][4]int) bool {

	//check duplicates by row
	for row := 0; row < 9; row++ {
		counter := [10]int{}
		for col := 0; col < 9; col++ {
			counter[board[row][col]]++
		}
		if hasDuplicates(counter) {
			return false
		}
	}

	//check duplicates by column
	for row := 0; row < 9; row++ {
		counter := [10]int{}
		for col := 0; col < 9; col++ {
			counter[board[col][row]]++
		}
		if hasDuplicates(counter) {
			return false
		}
	}

	//check 3x3 section
	for i := 0; i < 9; i += 3 {
		for j := 0; j < 9; j += 3 {
			counter := [10]int{}
			for row := i; row < i+3; row++ {
				for col := j; col < j+3; col++ {
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

func hasDuplicates([10]int) bool {

	return true
}

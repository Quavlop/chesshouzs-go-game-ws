package helpers

import "strings"

func ConvertNotationToArray(input string) [][]rune {
	rows := strings.Split(input, "|")
	result := make([][]rune, len(rows))

	for i, row := range rows {
		result[i] = []rune(row)
	}

	return result
}

func ConvertArrayToNotation(array [][]rune) string {
	var sb strings.Builder

	for _, row := range array {
		sb.WriteString(string(row))
		sb.WriteString("|")
	}

	if sb.Len() > 0 {
		return sb.String()[:sb.Len()-1]
	}
	return ""
}

func TransformBoard(array [][]rune) [][]rune {
	len := len(array)
	newState := make([][]rune, len)

	for i := range array {
		newState[i] = make([]rune, len)
		copy(newState[i], array[i])
	}

	for row := 0; row < len/2; row++ {
		for col := 0; col < len; col++ {
			temp := newState[row][col]
			newState[row][col] = newState[len-row-1][len-col-1]
			newState[len-row-1][len-col-1] = temp
		}
	}

	return newState
}

package cmd

func mapToKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
}

func averageStringLen(strings []string) int {
	var totalLen int = 0
	var cnt int = 0

	for _, str := range strings {
		currentLen := len([]rune(dropAnsiCodes(str)))
		totalLen += currentLen
		cnt += 1
	}

	if cnt == 0 {
		cnt = 1
	}

	return totalLen / cnt
}

func averageLineLen(lines []string) int {
	linesLen := len(lines)
	if linesLen > 1 {
		lines = lines[:linesLen-1] //Drop last line, as it might skew up average length
	}

	return averageStringLen(lines)
}

func deleteLastChar(input []rune) []rune {
	len := len(input)
	if len == 0 {
		return input
	} else {
		return input[:len-1]
	}
}

func containsChar(input []rune, char rune) bool {
	for _, item := range input {
		if item == char {
			return true
		}
	}
	return false
}

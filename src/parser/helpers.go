package parser

import (
	"strconv"
)

var hex_symbols = map[string]int{
	"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	"A": 10, "B": 11, "C": 12, "D": 13, "E": 14, "F": 15,
}

// Splits the array at a specific target
func FirstInstance(data []byte, target string) ([]byte, []byte) {

	index := -1
	found := 0

	for i := 0; i < len(data); i++ {

		if data[i] == target[found%len(target)] {
			found++

			if found == len(target) {
				index = i + 1
				break
			}

		} else {
			found = 0
		}
	}

	if index == -1 {
		return nil, nil
	}

	return data[0 : index-len(target)], data[index:]
}

func urldecode(s string) string {

	offset := 0

	for index, char := range s {
		if char != '%' {
			continue
		}

		i := index - offset

		hex := s[i+1 : i+3]

		decimal, err := strconv.ParseInt(hex, 16, 32)

		if err != nil {
			continue
		}

		seg1 := s[0:i]
		seg2 := s[i+3:]

		offset += 2

		s = seg1 + string(decimal) + seg2
	}

	return s
}

//first index: starts directly at that index
//last index: starts at that index minus 1

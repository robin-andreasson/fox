package parser

import "net/url"

//Splits the array at a specific target
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

func decode(s string) string {
	result, _ := url.QueryUnescape(s)

	return result
}

package parser

//Splits the first instance and returns
func FirstInstance(data []byte, _target string) ([]byte, []byte) {

	target := string(_target)
	index := -1
	found := 0

	for i := 0; i < len(data); i++ {

		if data[i] == target[found%len(target)] {
			found += 1

			if found == len(target) {
				index = i + 1
				break
			}

		} else {
			found = 0
		}
	}

	var segment_1 []byte
	var segment_2 []byte

	if index == -1 {
		segment_1 = data
	} else {
		segment_1 = data[0 : index-len(target)]
		segment_2 = data[index:]
	}

	return segment_1, segment_2
}

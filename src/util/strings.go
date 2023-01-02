package util

/**
 * Returns an int representing how different two strings are
 * The lower the number, the more similar the strings are
 */
func stringSimilarity(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	var cost int
	if a[len(a)-1] == b[len(b)-1] {
		cost = 0
	} else {
		cost = 1
	}

	res := min(
		stringSimilarity(a[:len(a)-1], b)+1,
		stringSimilarity(a, b[:len(b)-1])+1,
		stringSimilarity(a[:len(a)-1], b[:len(b)-1])+cost,
	)

	return res
}

func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

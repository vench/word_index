package word_index

func mergeOrderedArrayOr(a ...[]ItemID) []ItemID {
	maxLen := 0
	maxValue := ItemID(0)

	for j := 0; j < len(a); j++ {
		if len(a[j]) == 0 {
			continue
		}

		if len(a[j]) > maxLen {
			maxLen = len(a[j])
		}

		if maxValue < a[j][len(a[j])-1] {
			maxValue = a[j][len(a[j])-1]
		}
	}

	maxValue++

	offsets := make([]int, len(a))
	b := make([]ItemID, 0, maxLen)
	lastIndex := ItemID(-1)
	minValue := maxValue

	for {
		minIndexResult := -1

		for j := 0; j < len(a); j++ {
			if len(a[j]) <= offsets[j] {
				continue
			}

			if a[j][offsets[j]] < minValue {
				minValue = a[j][offsets[j]]
				minIndexResult = j
			}
		}

		if minIndexResult == -1 {
			break
		}

		if lastIndex < minValue {
			b = append(b, minValue)
			lastIndex = minValue
		}

		minValue = maxValue
		offsets[minIndexResult]++
	}
	return b
}

func mergeOrderedArrayAnd(a ...[]ItemID) []ItemID {
	b := make([]ItemID, 0)
	minIndex := 0

	for i := 1; i < len(a); i++ {
		if len(a[minIndex]) > len(a[i]) {
			minIndex = i
		}
	}

	offsets := make([]int, len(a))

	for _, v := range a[minIndex] {
		has := true
		for j := 0; j < len(a); j++ {
			if j == minIndex {
				continue
			}

			for ; offsets[j] < len(a[j]); offsets[j]++ {
				if a[j][offsets[j]] > v {
					has = false
					break
				}

				if has = a[j][offsets[j]] == v; has {
					break
				}
			}

			if !has {
				break
			}
		}

		if has {
			b = append(b, v)
		}
	}

	return b
}

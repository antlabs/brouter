package brouter

func (t *treeNode) sortChildren() {
	for i := 1; i < len(t.children); i++ {
		max := i
		for j := i + 1; j < len(t.children); j++ {
			if t.children[i].childNum < t.children[j].childNum {
				max = j
			}
		}

		if max != i {
			tmpNode := t.children[i]
			t.children[i] = t.children[max]
			t.children[max] = tmpNode

			tmpChar := t.charIndex[i]
			t.charIndex[i] = t.charIndex[max]
			t.charIndex[max] = tmpChar
		}
	}
}

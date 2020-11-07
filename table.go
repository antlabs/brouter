// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

type table struct {
	recogOffsetMap map[byte]int
	offsetToChar   map[int]byte
	recogOffset    [256]int
	pos            int
}

func (t *table) init() {
	if t.recogOffsetMap == nil {
		t.recogOffsetMap = make(map[byte]int)
	}

	if t.offsetToChar == nil {
		t.offsetToChar = make(map[int]byte)
	}

	t.getCodeOffsetAndInsert('/')
}

func (t *table) getCodeOffsetAndInsert(c byte) int {
	if offset, ok := t.recogOffsetMap[c]; ok {
		return offset
	}

	t.pos++
	t.recogOffsetMap[c] = t.pos
	t.recogOffset[c] = t.pos
	t.offsetToChar[t.pos] = c

	return t.pos
}

func (t *table) getOffsetToChar(offset int) byte {
	return t.offsetToChar[offset]
}

func (t *table) getCodeOffset(c byte) (offset int, found bool) {
	offset = t.recogOffset[c]
	return offset, offset != 0
}

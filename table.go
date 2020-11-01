// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

type table struct {
	recogOffsetMap map[byte]int
	recogOffset    [256]int
	pos            int
}

func (t *table) init() {
	if t.recogOffsetMap == nil {
		t.recogOffsetMap = make(map[byte]int)
	}

	t.pos = 1
}

func (t *table) getCodeOffsetInsert(c byte) int {
	if offset, ok := t.recogOffsetMap[c]; ok {
		return offset
	}

	t.pos++
	t.recogOffsetMap[c] = t.pos
	t.recogOffset[c] = t.pos

	return t.pos
}

// 返回0表示没有找到
func (t *table) getCodeOffsetLookup(c byte) int {
	return t.recogOffset[c]
}

var defaultTable table

func init() {
	defaultTable.init()
}

func getCodeOffsetInsert(c byte) int {
	return defaultTable.getCodeOffsetInsert(c)
}

func getCodeOffsetLookup(c byte) int {
	return defaultTable.getCodeOffsetLookup(c)
}

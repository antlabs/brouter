package brouter

type nodeType int

const (
	ordinary nodeType = iota
	param
	wildcard
)

type tree struct {
	root *treeNode
}

func (r *tree) insert(path string, h HandleFunc) {
	r.root.insert(path, h)
}

func (r *tree) lookup(path string, p *Params) {
	r.root.lookup(path, p)
}

type treeNode struct {
	children []*treeNode
	segment
}

func (n *treeNode) getChildrenInsert(c byte) (int, *treeNode) {
	offset := getCodeOffsetInsert(c)
	if offset >= len(n.children) {
		newChildren := make([]*treeNode, offset)
		copy(newChildren, n.children)
		n.children = newChildren
	}

	return offset, n.children[offset]
}

func (n *treeNode) noConflict() {
}

// 这里分几个状态
// 找到空的位置可以插入
// 有共同前缀需要分裂节点
func (n *treeNode) insert(path string, h HandleFunc) {
	p := genPath(path, h)

	for _, segment := range p.segments {

		if len(n.segment.path) == 0 {
			n.noConflict()
		}

		if segment.nodeType == param || segment.nodeType == wildcard {
			if n.paramOrWildcard != nil {
				panic("TODO 1, 检查paramName是否一样")
			}

			n.paramOrWildcard = &treeNode{
				segment: segment,
			}

			n = n.paramOrWildcard
			continue
		}

		offset, children := n.getChildrenInsert(segment.path[0])
		if children == nil {
			n.children[offset] = &treeNode{
				segment: segment,
			}

			n = n.children[offset]
			continue
		}

		//分裂, 找到共同前缀
		tailPath := n.path
		insertPath := segment.path
		i, j := 0, 0
		for i < len(tailPath) && j < len(insertPath) {
			if tailPath[i] != insertPath[j] {
				break
			}
			i++
			j++
		}
	}
}

func (n *treeNode) lookup(path string, p *Params) (h HandleFunc) {

	for i := 0; i < len(path); {
		// 当前节点path大于剩余需要匹配的path，说明路径和该节点不匹配
		if len(n.path) > len(path[i:]) {
			return nil
		}

		// 当前节点path和需要匹配的路径比较下，如果不相等，返回空handleFunc
		if n.path != path[i:len(n.path)] {
			return nil
		}

		i += len(n.path) // 跳过n.path的部分
		if len(n.children) != 0 && n.children[0] != nil {
			n = n.children[0]

			p.appendKey(n.paramName)

			if n.nodeType == param {
				j := i
				for ; j < len(path) && path[j] != '/'; j++ {
				}

				p.setVal(path[i:j])
				i = j
			}

			if n.nodeType == wildcard {
				p.setVal(path[i:len(path)])
				return n.handle
			}
		}

		c := path[i]
		offset := getCodeOffsetLookup(c)
		if offset >= len(n.children) {
			return nil
		}

		n = n.children[offset]
		if n == nil {
			return nil
		}
		i++

	}

	return n.handle
}

package brouter

type nodeType int

const (
	ordinary nodeType = iota
	param
	wildcard
)

/*
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

func (n *treeNode) insert(path string, h HandleFunc) {
	p := genPath(path, h)

	for _, segment := range p.segments {

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

func (n *treeNode) lookup(path string, p *Params) {

	for i := 0; i < len(path); i++ {
		if n.paramOrWildcard != nil {
			if n.paramOrWildcard.nodeType == param {
				j := i
				for ; j < len(path) && path[j] != '/'; j++ {
				}

			}

			if n.paramOrWildcard.nodeType == wildcard {
				p.appendKey(n.paramName)
				p.setVal(path[j:len(path)])
			}
		}

		c := path[i]
		offset := getCodeOffsetLookup(c)

	}
}
*/

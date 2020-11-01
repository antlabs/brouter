package brouter

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

func (n *treeNode) getParamOrWildcard() *treeNode {
	return n.getChildrenIndexAndMalloc(0)
}

func (n *treeNode) getChildrenNode(c byte) *treeNode {
	offset := getCodeOffsetInsert(c)
	return n.getChildrenIndexAndMalloc(offset)
}

func (n *treeNode) getChildrenIndexAndMalloc(index int) *treeNode {
	if index >= len(n.children) {
		newChildren := make([]*treeNode, index)
		copy(newChildren, n.children)
		n.children = newChildren
	}

	node := n.children[index]
	if node == nil {
		n.children[index] = &treeNode{}
	}
	return n.children[index]
}

func (n *treeNode) getNextTreeNode(i int, p path) (nextNode *treeNode) {
	if i+1 < len(p.segments) {
		nextSegment := p.segments[i+1]
		// 判断下个判断插入的节点类型
		// if 部分是param or wildcard 类型
		if nextSegment.nodeType.isParamOrWildcard() {
			nextNode = n.getParamOrWildcard()
		} else { // 这是普通节点
			nextNode = n.getChildrenNode(nextSegment.path[0])
		}
	}

	return
}

// 这里分几个状态
// 找到空的位置可以插入
// 有共同前缀需要分裂节点
func (n *treeNode) insert(path string, h HandleFunc) {
	p := genPath(path, h)

	for i := 0; i < len(p.segments); i++ {

		segment := p.segments[i]

		nextNode := n.getNextTreeNode(i, p)

		// 如果n.segment.path 为空，就可以直接插入到这个节点
		// 注意:区分普通节点和变量节点
		if n.nodeType.isEmpty() {
			// 普通节点
			if segment.nodeType.isOrdinary() {
				n.segment = segment
				n = nextNode
				continue
			}

			// 变量或者通配符节点
			if segment.nodeType.isParamOrWildcard() {
				n.path = segment.path
				paramOrWildcard := n.getParamOrWildcard()
				paramOrWildcard.segment = segment
				paramOrWildcard.path = ""
				n = paramOrWildcard.getNextTreeNode(i, p)
				continue
			}
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

		if i != len(tailPath) {
		}

		// TODO 处理变量节点
	}
}

func (n *treeNode) lookup(path string, p *Params) (h HandleFunc) {

	for i := 0; i < len(path); {
		// 当前节点path大于剩余需要匹配的path，说明路径和该节点不匹配
		if len(n.path) > len(path[i:]) {
			return nil
		}

		// 当前节点path和需要匹配的路径比较下，如果不相等，返回空指针
		if n.path != path[i:len(n.path)] {
			return nil
		}

		i += len(n.path)                                  // 跳过n.path的部分
		if len(n.children) != 0 && n.children[0] != nil { //检查参数部分
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

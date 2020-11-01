// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

// tree
type tree struct {
	root *treeNode
}

// 插入函数
func (r *tree) insert(path string, h HandleFunc) {
	r.root.insert(path, h)
}

// 查找函数
func (r *tree) lookup(path string, p *Params) {
	r.root.lookup(path, p)
}

// treeNode，查找树node
type treeNode struct {
	children []*treeNode
	segment
}

// 获取param or wildcard 节点
func (n *treeNode) getParamOrWildcard() *treeNode {
	return n.getChildrenIndexAndMalloc(0)
}

// 获取子节点
func (n *treeNode) getChildrenNode(c byte) *treeNode {
	offset := getCodeOffsetInsert(c)
	return n.getChildrenIndexAndMalloc(offset)
}

// 有子节点就返回，没有先分配空间然后返回指针
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
	if i < len(p.segments) {
		nextSegment := p.segments[i]
		// 判断下个判断插入的节点类型
		// 如果是param or wildcard 类型
		if nextSegment.nodeType.isParamOrWildcard() {
			nextNode = n.getParamOrWildcard()
		} else { // 这是普通节点
			nextNode = n.getChildrenNode(nextSegment.path[0])
		}
	}

	return
}

func (n *treeNode) directInsert(segment segment, nextNode *treeNode, i int, p path) (*treeNode, bool) {
	// 普通节点
	if segment.nodeType.isOrdinary() {
		n.segment = segment
		return nextNode, true
	}

	// 变量或者通配符节点
	if segment.nodeType.isParamOrWildcard() {
		n.path = segment.path
		paramOrWildcard := n.getParamOrWildcard()
		paramOrWildcard.segment = segment
		paramOrWildcard.path = ""
		return paramOrWildcard.getNextTreeNode(i, p), true
	}
	return nil, false
}

func (n *treeNode) splitNode(sm segment, i int, p path) *treeNode {
	//分裂, 找到共同前缀
	// 特殊情况已经被剔除掉，这里是需要新加子节点的情况
	tailPath := n.path
	insertPath := sm.path
	i, j := 0, 0
	for i < len(tailPath) && j < len(insertPath) {
		if tailPath[i] != insertPath[j] {
			break
		}
		i++
		j++
	}

	// 儿子变孙子
	grandson := n.children
	n.children = make([]*treeNode, 0, 1)

	splitSegment := segment{
		path:     n.path[i:],
		handle:   n.handle,
		nodeType: n.nodeType,
	}

	n.path = n.path[:i]
	n.handle = nil

	//TODO, 待插入segment如果是特殊节点，和n.path有重合的路径(n.path是普通路径), 直接panic

	if len(tailPath[i:]) > 0 {
		nextNode := n.getChildrenNode(tailPath[i])
		nextNode.children = grandson
		nextNode.segment = splitSegment
	} else {
		panic("splitNode:This is not taken into account")
	}

	if len(insertPath[j:]) > 0 {
		nextNode := n.getChildrenNode(tailPath[j])
		sm.path = insertPath[j:]
		nextNode.segment = sm
		return nextNode.getNextTreeNode(i+1, p)
	}

	n.handle = sm.handle

	return n.getNextTreeNode(i+1, p)
}

// 这里分情况讨论
func (n *treeNode) insert(path string, h HandleFunc) {
	p := genPath(path, h)

	for i := 0; i < len(p.segments); i++ {

		segment := p.segments[i]

		nextNode := n.getNextTreeNode(i+1, p)

		for {
			// 1.直接插入
			// 如果n.segment.isOrdinary() 为空，就可以直接插入到这个节点
			// 注意:区分普通节点和变量节点
			if n.nodeType.isEmpty() {
				if next, ok := n.directInsert(segment, nextNode, i+1, p); ok {
					nextNode = next
					break
				}

				panic("Unknown node type")
			}

			// 3,4,5 考虑下变量和可变参数
			// 3.不需要分裂 node, 当前需要插入的path和n.path相同
			if n.equal(segment) {
				if n.handle == nil {
					n.handle = segment.handle
				}
				n = nextNode
				break
			}

			// 4.不需要分裂 node, 当前需要插入的path包含n.path
			// 这种情况比较复杂, 剔除重复前缀元素，重走上面流程
			if len(n.path) < len(segment.path) && n.path == segment.path[:len(n.path)] {
				segment.path = segment.path[len(n.path):]
				n = n.getNextTreeNode(i, p)
				continue
			}

			// 5.分裂节点再插入
			n.splitNode(segment, i, p)
			break
		}

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

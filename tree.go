// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

import (
	"fmt"
	"strings"
	"sync"
)

// tree
type tree struct {
	root      *treeNode
	paramPool sync.Pool
	maxParam  int
}

// 构造一课树
func newTree() *tree {
	t := &tree{root: &treeNode{}}
	return t
}

// 插入函数
func (r *tree) insert(path string, h HandleFunc) {
	p := genPath(path, h)

	r.changePool(&p)
	r.root.insert(path, h, p)
}

// 获取参数
func (r *tree) getParam() *Params {
	ptr := r.paramPool.Get().(*Params)
	*ptr = (*ptr)[0:0]
	return ptr
}

// 查找函数
func (r *tree) lookup(path string) (HandleFunc, *Params) {
	return r.root.lookup(path, r.getParam)
}

// treeNode，查找树node
type treeNode struct {
	// 分段结构
	segment
	children []*treeNode
	// 当前节点有特殊节点的儿子, param or wildcard
	haveParamWildcardChild bool
	// char和 children索引的关系
	charIndex []byte
	// 儿子(直接联系的子节点)个数，不计算跨代的孩子节点，儿子多的排前面
	childNum int32
}

// 获取param or wildcard 节点
func (n *treeNode) getParamOrWildcard() *treeNode {
	return n.getChildrenIndexAndMalloc(0, 0)
}

func (n *treeNode) getChildrenNode2(char byte) (*treeNode, bool) {
	offset := -1
	for i, c := range n.charIndex {
		if c == char {
			offset = i
			break
		}
	}
	// 没有找到
	newNode := false
	if offset == -1 {
		newNode = true
		offset = len(n.charIndex)
		if offset == 0 {
			offset = 1
		}
	}

	return n.getChildrenIndexAndMalloc(offset, char), newNode
}

// 获取子节点
func (n *treeNode) getChildrenNode(char byte) *treeNode {
	newNode, _ := n.getChildrenNode2(char)
	return newNode
}

// 分配需要的节点, 有就返回，没有就先分配再返回
func (n *treeNode) getChildrenIndexAndMalloc(index int, char byte) *treeNode {
	if index >= len(n.children) {
		newChildren := make([]*treeNode, index+1)
		copy(newChildren, n.children)
		n.children = newChildren

	}

	if index >= len(n.charIndex) {
		newCharIndex := make([]byte, index+1)
		copy(newCharIndex, n.charIndex)
		n.charIndex = newCharIndex
	}

	node := n.children[index]
	if node == nil {
		n.children[index] = &treeNode{}
	}

	char2 := n.charIndex[index]
	if char2 == 0 {
		n.charIndex[index] = char
	}

	return n.children[index]
}

// 返回需要插入的儿子节点和下个需要插入的节点类型
func (n *treeNode) getNextTreeNodeAndType(i int, p path) (nextNode *treeNode, nextNodeType nodeType) {
	if i < len(p.segments) {
		nextSegment := p.segments[i]
		if len(nextSegment.path) > 0 {
			return n.getChildrenNode(nextSegment.path[0]), param
		}

		return n.getParamOrWildcard(), nextSegment.nodeType
	}

	return
}

// 返回需要插入的儿子节点
func (n *treeNode) getNextTreeNode(i int, p path) (nextNode *treeNode) {
	if i < len(p.segments) {
		nextSegment := p.segments[i]
		if len(nextSegment.path) > 0 {
			return n.getChildrenNode(nextSegment.path[0])
		}

		return n.getParamOrWildcard()
	}

	return
}

// 如果是空节点就直接插入
func (n *treeNode) directInsert(segment segment, i int, p path) (*treeNode, bool) {
	// 普通节点, 或者没有path的变量节点
	if segment.nodeType.isOrdinary() || len(segment.path) == 0 && segment.nodeType.isParamOrWildcard() {
		n.segment = segment

		nextNode, nextType := n.getNextTreeNodeAndType(i, p)
		if nextType.isParamOrWildcard() {
			n.haveParamWildcardChild = true
		}
		return nextNode, true
	}

	// 变量或者通配符节点
	if segment.nodeType.isParamOrWildcard() {
		// TODO segment.path 为空的param or wildcard
		//if n.path != "" {
		n.path = segment.path
		n.nodeType = ordinary

		n.haveParamWildcardChild = true
		paramOrWildcard := n.getParamOrWildcard()
		paramOrWildcard.childNum++
		paramOrWildcard.segment = segment
		paramOrWildcard.path = ""

		return paramOrWildcard.getNextTreeNode(i, p), true
		//}

	}

	return nil, false
}

// 分裂当前节点
func (n *treeNode) splitCurrentNode(i int) {
	// 儿子变孙子
	grandson := n.children
	grandsonCharIndex := n.charIndex
	grandsomChildNum := n.childNum - 1
	n.children = make([]*treeNode, 0, 1)
	n.charIndex = make([]byte, 0, 1)

	splitSegment := segment{
		path:     n.path[i:],
		handle:   n.handle,
		nodeType: n.nodeType,
	}

	c := n.path[i]
	n.path = n.path[:i]
	n.handle = nil

	nextNode := n.getChildrenNode(c)
	nextNode.children = grandson
	nextNode.charIndex = grandsonCharIndex
	nextNode.childNum = grandsomChildNum
	nextNode.segment = splitSegment
	nextNode.haveParamWildcardChild = n.haveParamWildcardChild

	n.haveParamWildcardChild = false
}

func prefixIndex(s1, s2 string) int {
	i := 0
	for i < len(s1) && i < len(s2) {
		if s1[i] != s2[i] {
			return i
		}
		i++
	}
	return i
}

func (n *treeNode) splitNode(sm *segment, segIndex int, p path) (*treeNode, bool) {
	//分裂, 找到共同前缀
	// 特殊情况已经被剔除掉，这里是需要新加子节点的情况
	insertPath := sm.path

	// 找到相同前缀下标
	i := prefixIndex(n.path, insertPath)

	//fmt.Printf("n.path:%s, insertPath:%s, handle:%p\n", n.path, sm.path, n.handle)
	//TODO, 待插入segment如果是特殊节点，和n.path有重合的路径(n.path是普通路径), 直接panic
	// 1. n 是特殊节点，insertPath是普通节点
	// 2. n 是普通节点，insertPath是特殊节点

	if len(n.path[i:]) > 0 {
		n.splitCurrentNode(i)
	} // else n.path只有'/'

	sm.path = insertPath[i:]
	if len(insertPath[i:]) > 0 {
		nextNode, newNode := n.getChildrenNode2(insertPath[i])
		if newNode {
			nextNode.childNum++
		}
		return nextNode, false
	}

	// insertPath 为0的情况
	n.handle = sm.handle
	if sm.nodeType.isParamOrWildcard() {
		n.haveParamWildcardChild = true
		paramOrWildcard := n.getParamOrWildcard()

		if paramOrWildcard.paramName != "" && paramOrWildcard.paramName != sm.paramName {
			panic(fmt.Sprintf("paramName: %s", sm.paramName))
		}

		if paramOrWildcard.handle == nil {
			paramOrWildcard.handle = sm.handle
		}

		paramOrWildcard.path = ""
		return paramOrWildcard.getNextTreeNode(segIndex, p), true
	}

	// /repos/:owner/:repo/milestones/:number/labels
	// /repos/:owner/:repo/milestones
	return n.getNextTreeNode(segIndex, p), true
}

// 这里分情况讨论
func (n *treeNode) insert(path string, h HandleFunc, p path) {

	if h2, _ := n.lookup(path, nil); h2 != nil {
		//panic("duplicate registration '" + path + "'")
	}

	for i := 0; i < len(p.segments); i++ {

		segment := p.segments[i]

		prevNode := n
		n.childNum++
		fmt.Printf("for node:%p\n", n)
		for {
			// 1.直接插入
			// 如果n.segment.isOrdinary() 为空，就可以直接插入到这个节点
			// 注意:区分普通节点和变量节点
			if n.nodeType.isEmpty() {
				fmt.Printf("0. isEmpty:  [%8s(:%s)]:[%s], [a:%p], [cn:%d]\n",
					segment.path, segment.paramName, path, n, n.childNum)
				if next, ok := n.directInsert(segment, i+1, p); ok {
					if next != nil {
						//next.childNum++
					}
					n = next
					break
				}

				panic(fmt.Sprintf("Unknown node type:%s", segment.nodeType))
			}

			// 2,3,4 考虑下变量和可变参数
			// 2.不需要分裂 node, 当前需要插入的path和n.path相同
			if n.equal(segment) {
				fmt.Printf("2. equal:    [%8s(:%s)]:[%s], [a:%p], [cn:%d]\n",
					segment.path, segment.paramName, path, n, n.childNum)
				if n.handle == nil {
					n.handle = segment.handle
				}
				if prevNode != n {
					n.childNum++
				}
				n = n.getNextTreeNode(i+1, p)
				break
			}

			// 3.不需要分裂 node, 当前需要插入的path包含n.path
			// 这种情况比较复杂, 剔除重复前缀元素，重走上面流程
			if len(n.path) < len(segment.path) && n.path == segment.path[:len(n.path)] {

				fmt.Printf("3. contain:  [%8s(:%s)]:[%s], [a:%p], [cn:%d], n.path:%s, segment.path%s\n",
					segment.path, segment.paramName, path, n, n.childNum, n.path, segment.path)

				segment.path = segment.path[len(n.path):]
				n2, newNode := n.getChildrenNode2(segment.path[0])
				if prevNode != n {
					n.childNum++
				}
				if newNode {
					n2.childNum++
				}
				n = n2
				continue
			}

			// 4.分裂节点再插入
			fmt.Printf("4.splitNode: [%8s(:%s)]:[%s], [a:%p] [cn:%d]\n",
				segment.path, segment.paramName, path, n, n.childNum)
			var next bool
			if n, next = n.splitNode(&segment, i+1, p); next {
				break
			}
		}

		prevNode.sortChildren()
	}
}

func (n *treeNode) lookup(path string, getParam func() *Params) (h HandleFunc, p *Params) {

next:
	for {

		// 当前节点path大于剩余需要匹配的path，说明路径和该节点不匹配
		if len(n.path) > len(path) {
			return nil, p
		}

		// 当前节点path和需要匹配的路径比较下，如果不相等，返回空指针
		if n.path != path[:len(n.path)] {
			return nil, p
		}

		// 长度相等
		if len(path) == len(n.path) {
			return n.handle, p
		}

		path = path[len(n.path):]

		// 普通节点
		if !n.haveParamWildcardChild {
			for i, c := range n.charIndex {
				if c == path[0] {
					n = n.children[i]
					continue next
				}

			}

			return n.handle, p

		}

		// 特殊节点
		n = n.children[0]

		if p == nil && getParam != nil {
			p = getParam()
		}

		if n.nodeType == param {
			j := 0
			for ; j < len(path) && path[j] != '/'; j++ {
			}

			if p != nil {
				*p = append(*p, Param{Key: n.paramName, Value: path[:j]})
			}

			if j == len(path) {
				return n.handle, p
			}

			if len(n.children) < 2 {
				return nil, p
			}

			path = path[j:]
			// TODO 排序可能会有冲突
			n = n.children[1] // '/'

			continue
		}

		if n.nodeType == wildcard {
			if p != nil {
				*p = append(*p, Param{Key: n.paramName, Value: path})
			}
			return n.handle, p
		}

	}

	return n.handle, p
}

func (t *tree) changePool(p *path) {
	if t.paramPool.New == nil {
		t.paramPool.New = func() interface{} {
			p := make(Params, 0, 0)
			return &p
		}
	}

	if p.maxParam > t.maxParam {
		t.maxParam = p.maxParam
		t.paramPool.New = func() interface{} {
			p := make(Params, 0, t.maxParam)
			return &p
		}
	}

}

// 打印一个节点
func (n *treeNode) info() {
	fmt.Printf("\n ============== start treeNode ######, %p\n", n)
	fmt.Printf("	n.path:%s\n", n.path)
	fmt.Printf("	children: ")
	for i := 0; i < len(n.children); i++ {
		fmt.Printf("%p ", n.children[i])
	}
	fmt.Printf("\n")

	fmt.Printf("	char    : [")
	for i := 0; i < len(n.charIndex); i++ {
		c := n.charIndex[i]
		fmt.Printf("%c, ", c)
	}
	fmt.Printf("]\n")

	fmt.Printf("	children-number: [")
	for i := 0; i < len(n.charIndex); i++ {
		num := int32(0)
		if n.children[i] != nil {
			num = n.children[i].childNum
		}

		fmt.Printf("%d, ", num)
	}
	fmt.Printf("]\n")

	fmt.Printf("	charIndex:%s\n", n.charIndex)
	fmt.Printf(" ==============   end treeNode ######, %p\n\n", n)
}

// 有效子节点个数
func (n *treeNode) childrenLen() (count int) {
	for _, n1 := range n.children {
		if n1 == nil {
			continue
		}
		count++
	}
	return
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// 打印当前节点的子树
// 有子节点用蓝色表示
// 变量节点使用红色表示
func (n *treeNode) depthfirst(depth int, prev int) {
	// 尾巴节点的兄弟节点
	prefix := strings.Repeat(" ", depth+prev)

	// 构造节点名字
	nodeName := n.path
	if n.paramName != "" {
		nodeName = ":" + n.paramName
	}

	//nodeName += fmt.Sprintf("[n:%d]", n.childNum)
	addr := fmt.Sprintf("%p", n)
	nodeName += fmt.Sprintf("[n:%d][a:%s]", n.childNum, addr[len(addr)-4:len(addr)])
	// 打印节点名字
	fmt.Println(prefix + nodeName)

	for _, n1 := range n.children {
		if n1 == nil {
			continue
		}

		n1.depthfirst(depth+1, prev+len(nodeName))

	}
}

// 打印一颗数
func (t *tree) info() {
	t.root.depthfirst(0, 0)
}

// 查询节点, 测试用
func (n *treeNode) lookupNode(path string) *treeNode {

next:
	for {

		// 当前节点path大于剩余需要匹配的path，说明路径和该节点不匹配
		if len(n.path) > len(path) {
			return nil
		}

		// 当前节点path和需要匹配的路径比较下，如果不相等，返回空指针
		if n.path != path[:len(n.path)] {
			return nil
		}

		// 长度相等
		if len(path) == len(n.path) {
			return n
		}

		path = path[len(n.path):]

		// 普通节点
		if !n.haveParamWildcardChild {
			for i, c := range n.charIndex {
				if c == path[0] {
					n = n.children[i]
					continue next
				}

			}

			return n

		}

		// 特殊节点
		n = n.children[0]

		if n.nodeType == param {
			j := 0
			for ; j < len(path) && path[j] != '/'; j++ {
			}

			if j == len(path) {
				return n
			}

			if len(n.children) < 2 {
				return nil
			}

			path = path[j:]
			// TODO 排序可能会有冲突
			n = n.children[1] // '/'

			continue
		}

		if n.nodeType == wildcard {
			return n
		}

	}

	return n
}

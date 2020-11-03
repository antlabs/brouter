// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

import (
	"fmt"
	"sync"
)

// tree
type tree struct {
	root      *treeNode
	paramPool sync.Pool
	maxParam  int
}

func newTree() *tree {
	return &tree{root: &treeNode{}}
}

// 插入函数
func (r *tree) insert(path string, h HandleFunc) {
	p := genPath(path, h)
	r.changePool(&p)
	r.root.insert(path, h, p)
}

// 查找函数
func (r *tree) lookup(path string, p *Params) HandleFunc {
	return r.root.lookup(path, p)
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
		newChildren := make([]*treeNode, index+1)
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
		nextNode = n.getChildrenNode(nextSegment.path[0])
	}

	return
}

func (n *treeNode) directInsert(segment segment, i int, p path) (*treeNode, bool) {
	// 普通节点
	if segment.nodeType.isOrdinary() {
		n.segment = segment

		return n.getNextTreeNode(i+1, p), true
	}

	// 变量或者通配符节点
	if segment.nodeType.isParamOrWildcard() {
		n.path = segment.path
		n.nodeType = ordinary

		paramOrWildcard := n.getParamOrWildcard()
		paramOrWildcard.segment = segment
		paramOrWildcard.path = ""

		return paramOrWildcard.getNextTreeNode(i, p), true
	}

	return nil, false
}

func (n *treeNode) splitNode(sm segment, segIndex int, p path) *treeNode {
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
	// 1. n 是特殊节点，insertPath是普通节点
	// 2. n 是普通节点，insertPath是特殊节点

	if len(tailPath[i:]) > 0 {
		nextNode := n.getChildrenNode(tailPath[i])
		nextNode.children = grandson
		nextNode.segment = splitSegment
	} else {
		panic(fmt.Sprintf("splitNode:This is not taken into account:tailPath:%s", tailPath))
	}

	if len(insertPath[j:]) > 0 {
		nextNode := n.getChildrenNode(insertPath[j])
		sm.path = insertPath[j:]
		nextNode.segment = sm
		return nextNode.getNextTreeNode(segIndex, p)
	}

	// insertPath 为0的情况
	n.handle = sm.handle

	return n.getNextTreeNode(segIndex, p)
}

// 这里分情况讨论
func (n *treeNode) insert(path string, h HandleFunc, p path) {

	for i := 0; i < len(p.segments); i++ {

		segment := p.segments[i]

		for {
			// 1.直接插入
			// 如果n.segment.isOrdinary() 为空，就可以直接插入到这个节点
			// 注意:区分普通节点和变量节点
			if n.nodeType.isEmpty() {
				if next, ok := n.directInsert(segment, i+1, p); ok {
					n = next
					break
				}

				panic("Unknown node type")
			}

			// 2,3,4 考虑下变量和可变参数
			// 2.不需要分裂 node, 当前需要插入的path和n.path相同
			if n.equal(segment) {
				if n.handle == nil {
					n.handle = segment.handle
				}
				n = n.getNextTreeNode(i+1, p)
				break
			}

			// 3.不需要分裂 node, 当前需要插入的path包含n.path
			// 这种情况比较复杂, 剔除重复前缀元素，重走上面流程
			if len(n.path) < len(segment.path) && n.path == segment.path[:len(n.path)] {
				segment.path = segment.path[len(n.path):]
				n = n.getNextTreeNode(i, p)
				continue
			}

			// 4.分裂节点再插入

			n = n.splitNode(segment, i+1, p)
			break
		}

	}
}

// 判断变量节点
func (n *treeNode) hasParamOrWildcard() bool {
	return len(n.children) != 0 &&
		n.children[0] != nil &&
		n.children[0].nodeType.isParamOrWildcard()
}

func (n *treeNode) getChildrenIndex(c byte) *treeNode {
	offset := getCodeOffsetLookup(c)
	if offset >= len(n.children) {
		return nil
	}

	return n.children[offset]

}

func (n *treeNode) debug() {
	fmt.Printf("\n %p, ============== start treeNode ######\n", n)
	fmt.Printf("	n.path:%s\n", n.path)
	fmt.Printf("	children: ")
	for i := 0; i < len(n.children); i++ {
		fmt.Printf("%p ", n.children[i])
	}
	fmt.Printf("\n")

	fmt.Printf("	char    : [")
	for i := 0; i < len(n.children); i++ {
		fmt.Printf("%c, ", getOffsetToChar(i))
	}
	fmt.Printf("]\n")
	fmt.Printf(" %p, ==============   end treeNode ######\n", n)
}

func (n *treeNode) checkPath(j *int, path string) (h HandleFunc, quit bool) {
	i := *j
	// 当前节点path大于剩余需要匹配的path，说明路径和该节点不匹配
	if len(n.path) > len(path[i:]) {
		return nil, true
	}

	// 当前节点path和需要匹配的路径比较下，如果不相等，返回空指针
	if n.path != path[i:i+len(n.path)] {
		return nil, true
	}

	*j += len(n.path) // 跳过n.path的部分

	if *j == len(path) {
		return n.handle, true
	}

	return nil, false
}

func (n *treeNode) lookup(path string, p *Params) (h HandleFunc) {

	for i := 0; i < len(path); {

		n.debug()

		if n.hasParamOrWildcard() {
			if h, quit := n.checkPath(&i, path); quit {
				return h
			}
			n = n.children[0]

			p.appendKey(n.paramName)

			if n.nodeType == param {
				j := i
				for ; j < len(path) && path[j] != '/'; j++ {
				}

				p.setVal(path[i:j])
				i = j

				if j == len(path) {
					return n.handle
				}

				n = n.getChildrenIndex(path[i])
				if n == nil {
					return nil
				}
				continue
			}

			if n.nodeType == wildcard {
				p.setVal(path[i:len(path)])
				return n.handle
			}
		}

		if h, quit := n.checkPath(&i, path); quit {
			return h
		}

		c := path[i]
		if n = n.getChildrenIndex(c); n == nil {
			return nil
		}

		//i++

	}

	return n.handle
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

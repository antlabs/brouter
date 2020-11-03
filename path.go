// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

import (
	"fmt"
	"strings"
)

type segment struct {
	path      string
	nodeType  nodeType
	paramName string
	handle    HandleFunc
}

func (s *segment) equal(s1 segment) bool {
	return s.path == s1.path && s.nodeType == s1.nodeType && s.paramName == s.paramName

}

type path struct {
	segments []segment
	maxParam int
}

// 标准库自带的split会去除分割符，但是这里需要，所以写个自定义的split
func splitPath(p string) []string {
	if p[0] != '/' {
		panic("The path must start with /")
	}

	path := make([]string, 0, 2)
	prevIndex := 0
	for i := 0; i < len(p); i++ {
		c := p[i]
		if i != 0 && c == '/' {

			path = append(path, p[prevIndex:i])
			prevIndex = i
		}
	}

	if prevIndex != len(p) {
		path = append(path, p[prevIndex:len(p)])
	}

	return path
}

// 基于插入是比较少的场景，所以genPath函数没有做性能优化
func genPath(p string, h HandleFunc) (path path) {

	paths := splitPath(p)
	prevIndex := 0

	genPrevPath := func(paths []string, prevIndex, i int) string {
		var prevPath strings.Builder
		for _, p := range paths[prevIndex:i] {
			prevPath.WriteString(p)
		}
		return prevPath.String()
	}

	for i := 0; i < len(paths); i++ {

		v := paths[i]
		if len(v) == 0 {
			panic("TODO panic")
		}

		if v[1] == ':' || v[1] == '*' {
			if len(v) == 1 {
				panic(fmt.Sprintf("Parameter cannot be empty path:%s", p))
			}

			if prevIndex != i {
				prevPath := genPrevPath(paths, prevIndex, i)

				path.segments = append(path.segments, segment{path: prevPath, nodeType: ordinary})
			}

			nType := param
			if v[0] == '*' {
				nType = wildcard
				if len(paths[i:]) > 0 {
					panic(fmt.Sprintf("The wildcard symbol cannot be the middle element of the path:%s", p))
				}
			}

			if len(v) <= 2 {
				panic("Wrong parameter name")
			}

			path.segments = append(path.segments, segment{path: "/", nodeType: nType, paramName: v[2:len(v)]})

			prevIndex = i + 1
			path.maxParam++
		}
	}

	if prevIndex != len(paths) {
		prevPath := genPrevPath(paths, prevIndex, len(paths))
		path.segments = append(path.segments, segment{path: prevPath, nodeType: ordinary})
	}

	path.segments[len(path.segments)-1].handle = h
	return
}

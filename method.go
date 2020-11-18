package brouter

import (
	"errors"
)

var ErrMethod = errors.New("error method")

type method struct {
	method [7]*tree
}

func (m *method) init() {
	for k := range m.method {
		if m.method[k] == nil {
			m.method[k] = newTree()
		}
	}
}

func methodIndex(method string) (int, error) {
	if len(method) <= 2 {
		return 0, ErrMethod
	}

	switch method[0] {
	case 'G':
		return 0, nil
	case 'P':

		switch method[1] {
		case 'O':
			return 1, nil
		case 'U':
			return 2, nil
		case 'A':
			return 3, nil
		default:
			return 0, ErrMethod
		}
	case 'D':
		return 4, nil
	case 'H':
		return 5, nil
	case 'O':
		return 6, nil
	default:
		return 0, ErrMethod
	}

	return 0, ErrMethod
}

func (m *method) getTree(method string) *tree {
	index, err := methodIndex(method)
	if err != nil {
		return nil
	}

	return m.method[index]
}

func (m *method) save(method, path string, h HandleFunc) {
	index, err := methodIndex(method)
	if err != nil {
		panic(err.Error())
	}

	m.method[index].insert(path, h)
}

/*
type method struct {
	method map[string]*tree
}

func (m *method) init() {
	if m.method == nil {
		m.method = make(map[string]*tree)
	}

	for _, k := range []string{"GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT", "PATCH"} {
		if m.method[k] == nil {
			m.method[k] = newTree()
		}
	}
}

func (m *method) getTree(method string) *tree {
	t, _ := m.method[method]
	return t
}

func (m *method) save(method, path string, h HandleFunc) {
	m.method[method].insert(path, h)
}
*/

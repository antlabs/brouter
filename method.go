package brouter

import (
	"errors"
	"fmt"
)

type method struct {
	method [7]*tree
}

var ErrMethod = errors.New("error method")

func (m *method) init() {
	for k := range m.method {
		if m.method[k] == nil {
			m.method[k] = newTree()
		}
	}
}

func methodIndex(method string) (int, error) {
	if len(method) == 0 {
		return 0, ErrMethod
	}

	switch method[0] {
	case 'G':
		return 0, nil
	case 'P':
		if len(method) <= 1 {
			return 0, fmt.Errorf("%w, %s", ErrMethod, method)
		}

		switch method[1] {
		case 'O':
			return 1, nil
		case 'U':
			return 2, nil
		case 'A':
			return 3, nil
		default:
			return 0, fmt.Errorf("%w, %s", ErrMethod, method)
		}
	case 'D':
		return 4, nil
	case 'H':
		return 5, nil
	case 'O':
		return 6, nil
	default:
		return 0, fmt.Errorf("%w, %s", ErrMethod, method)
	}

	return 0, fmt.Errorf("%w, %s", ErrMethod, method)
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

// Apache-2.0 License
// Copyright [2020] [guonaihong]

package brouter

type nodeType int

const (
	empty nodeType = iota
	ordinary
	param
	wildcard
)

func (n nodeType) isEmpty() bool {
	return n == empty
}

func (n nodeType) isOrdinary() bool {
	return n == ordinary
}

func (n nodeType) isParamOrWildcard() bool {
	return n == param || n == wildcard
}

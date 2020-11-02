// Apache-2.0 License
// Copyright [2020] [guonaihong]
package brouter

import "net/http"

type Param struct {
	Key   string
	Value string
}

type Params []Param

func (ps Params) ByName(name string) string {
	for _, p := range ps {
		if p.Key == name {
			return p.Value
		}
	}
	return ""
}

// 深度拷贝
func (ps Params) Clone() Params {
	rv := make(Params, len(ps))
	for i := range ps {
		rv[i] = ps[i]
	}
	return rv
}

func (p *Params) appendKey(key string) {

	*p = append(*p, Param{Key: key})
}

func (p *Params) setVal(val string) {
	(*p)[len(*p)-1].Value = val
}

type HandleFunc func(w http.ResponseWriter, r *http.Request, p Params)

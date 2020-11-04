package brouter

import (
	"net/http"
)

type router struct {
	method
}

func New() *router {
	r := &router{}
	r.init()
	return r
}

// GET 请求
func (r *router) GET(path string, handle HandleFunc) {
	r.Handle(http.MethodGet, path, handle)
}

// HEAD 请求
func (r *router) HEAD(path string, handle HandleFunc) {
	r.Handle(http.MethodHead, path, handle)
}

// POST 请求
func (r *router) POST(path string, handle HandleFunc) {
	r.Handle(http.MethodPost, path, handle)
}

// PUT 请求
func (r *router) PUT(path string, handle HandleFunc) {
	r.Handle(http.MethodPut, path, handle)
}

// PATCH 请求
func (r *router) PATCH(path string, handle HandleFunc) {
	r.Handle(http.MethodPatch, path, handle)
}

// DELETE 请求
func (r *router) DELETE(path string, handle HandleFunc) {
	r.Handle(http.MethodDelete, path, handle)
}

// OPTIONS 请求
func (r *router) OPTIONS(path string, handle HandleFunc) {
	r.Handle(http.MethodOptions, path, handle)
}

func (r *router) Handle(method, path string, handle HandleFunc) {
	r.save(method, path, handle)
}

// 如果Params的生命周期超过ServeHTTP函数的，请clone一份Params
// 或者取走感兴趣的参数
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	tree := r.getTree(req.Method)
	if tree != nil {

		p2 := tree.paramPool.Get().(*Params)
		*p2 = (*p2)[0:0]
		put := func(p *Params) {
			tree.paramPool.Put(p)
		}

		handle := tree.lookup(path, nil)
		if handle != nil {
			handle(w, req, nil)
			put(p2)
			return
		}

		put(p2)
	}

	http.NotFound(w, req)
}

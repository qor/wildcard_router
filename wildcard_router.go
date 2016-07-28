package wildcard_router

import (
	"net/http"
)

type WildcardRouter struct {
	Handlers []WildcardInterface
}

type WildcardInterface interface {
	Handle(w http.ResponseWriter, req *http.Request) bool
}

func New(mux *http.ServeMux) *WildcardRouter {
	w := &WildcardRouter{}
	mux.Handle("/", w)
	return w
}

func (w *WildcardRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	w.WildcardHandle(writer, req)
}

func (w *WildcardRouter) AddHandler(handler WildcardInterface) {
	w.Handlers = append(w.Handlers, handler)
}

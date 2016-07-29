package wildcard_router

import (
	"net/http"
)

// WildcardRouter is holder of registed handlers
type WildcardRouter struct {
	Handlers []WildcardInterface
}

// WildcardInterface defined interfaces using to handle a router
type WildcardInterface interface {
	Handle(w http.ResponseWriter, req *http.Request) bool
}

// New will create a WildcardRouter and mount wildcard router to mux
func New(mux *http.ServeMux) *WildcardRouter {
	w := &WildcardRouter{}
	mux.Handle("/", w)
	return w
}

func (w *WildcardRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	w.WildcardHandle(writer, req)
}

// AddHandler will append a new handler to Handlers
func (w *WildcardRouter) AddHandler(handler WildcardInterface) {
	w.Handlers = append(w.Handlers, handler)
}

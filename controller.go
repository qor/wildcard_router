package wildcard_router

import (
	"net/http"
)

func (w *WildcardRouter) WildcardHandle(writer http.ResponseWriter, req *http.Request) {
	var matched bool
	for _, handler := range w.Handlers {
		if handler.Handle(writer, req) {
			matched = true
			break
		}
	}
	if !matched {
		http.NotFound(writer, req)
	}
}

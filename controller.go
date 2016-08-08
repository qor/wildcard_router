package wildcard_router

import "net/http"

// WildcardHandle will loop handlers to handle a request
func (w *WildcardRouter) WildcardHandle(writer *WildcardRouterWriter, req *http.Request) {
	var matched bool
	for _, handler := range w.Handlers {
		handler.ServeHTTP(writer, req)
		if writer.matchedStatus() {
			matched = true
			break
		}
		writer.reset()
	}
	if !matched {
		writer.ForceNotFound(req)
		writer.reset()
	}
}

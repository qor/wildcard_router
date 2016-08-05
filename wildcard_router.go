package wildcard_router

import "net/http"

// WildcardRouter is holder of registed handlers
type WildcardRouter struct {
	Handlers []http.Handler
}

type WildcardRouterWriter struct {
	http.ResponseWriter
	tmpStatus       int
	finalStatus     int
	skipCheckStatus bool
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
	w.WildcardHandle(&WildcardRouterWriter{writer, 0, 0, false}, req)
}

// AddHandler will append a new handler to Handlers
func (w *WildcardRouter) AddHandler(handler http.Handler) {
	w.Handlers = append(w.Handlers, handler)
}

func (w *WildcardRouterWriter) WriteHeader(statusCode int) {
	if w.skipCheckStatus || statusCode != http.StatusNotFound {
		w.skipCheckStatus = false
		w.finalStatus = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
	w.tmpStatus = statusCode
}

func (w *WildcardRouterWriter) Write(data []byte) (int, error) {
	if w.tmpStatus != http.StatusNotFound {
		return w.ResponseWriter.Write(data)
	} else {
		return 0, nil
	}
}

func (w WildcardRouterWriter) MatchedStatus() bool {
	return w.Status() != http.StatusNotFound && w.Status() != 0
}

func (w *WildcardRouterWriter) FocusNotFound(req *http.Request) {
	w.skipCheckStatus = true
	http.NotFound(w, req)
}

func (w WildcardRouterWriter) Status() int {
	return w.finalStatus
}

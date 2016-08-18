package wildcard_router

import "net/http"

// WildcardRouter holds registered route handlers
type WildcardRouter struct {
	Handlers []http.Handler
}

// New will create a WildcardRouter and mount it to http mux
func New(mux *http.ServeMux) *WildcardRouter {
	w := &WildcardRouter{}
	mux.Handle("/", w)
	return w
}

// AddHandler will append new handler to Handlers
func (w *WildcardRouter) AddHandler(handler http.Handler) {
	w.Handlers = append(w.Handlers, handler)
}

func (w *WildcardRouter) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	wildcardRouterWriter := &WildcardRouterWriter{writer, 0, false}

	for _, handler := range w.Handlers {
		if handler.ServeHTTP(wildcardRouterWriter, req); wildcardRouterWriter.isProcessed() {
			return
		}
		wildcardRouterWriter.reset()
	}

	wildcardRouterWriter.skipNotFoundCheck = true
	http.NotFound(wildcardRouterWriter, req)
}

// WildcardRouterWriter will used to capture status
type WildcardRouterWriter struct {
	http.ResponseWriter
	// Keep status code
	status int
	// Used to skip status check
	skipNotFoundCheck bool
}

// Status will return request's status code
func (w WildcardRouterWriter) Status() int {
	return w.status
}

// WriteHeader only set status code when not 404
func (w *WildcardRouterWriter) WriteHeader(statusCode int) {
	if w.skipNotFoundCheck || statusCode != http.StatusNotFound {
		w.ResponseWriter.WriteHeader(statusCode)
	}
	w.status = statusCode
}

// Write only set content when not 404
func (w *WildcardRouterWriter) Write(data []byte) (int, error) {
	if w.skipNotFoundCheck || w.status != http.StatusNotFound {
		w.status = http.StatusOK
		return w.ResponseWriter.Write(data)
	}
	return 0, nil
}

func (w *WildcardRouterWriter) reset() {
	w.skipNotFoundCheck = false
	w.status = 0
}

func (w WildcardRouterWriter) isProcessed() bool {
	return w.status != http.StatusNotFound && w.status != 0
}

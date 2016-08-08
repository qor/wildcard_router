package wildcard_router

import "net/http"

// WildcardRouter is holder of registed handlers
type WildcardRouter struct {
	Handlers []http.Handler
}

// WildcardRouterWriter will used to capture status
type WildcardRouterWriter struct {
	http.ResponseWriter
	// Storage status for each handler, will be reset after goes to next handler
	tmpStatus int
	// Storage the real status
	finalStatus int
	// Used to skip status check
	skipCheckStatus bool
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

// WriteHeader will only set status code if request isn't 404
func (w *WildcardRouterWriter) WriteHeader(statusCode int) {
	if w.skipCheckStatus || statusCode != http.StatusNotFound {
		w.finalStatus = statusCode
		w.ResponseWriter.WriteHeader(statusCode)
	}
	w.tmpStatus = statusCode
}

// Write will only set content if request isn't 404
func (w *WildcardRouterWriter) Write(data []byte) (int, error) {
	if w.skipCheckStatus || w.tmpStatus != http.StatusNotFound {
		w.finalStatus = http.StatusOK
		return w.ResponseWriter.Write(data)
	}
	return 0, nil
}

// ForceNotFound will force set request as not found
func (w *WildcardRouterWriter) ForceNotFound(req *http.Request) {
	w.skipCheckStatus = true
	http.NotFound(w, req)
}

// Status will return request's status code
func (w WildcardRouterWriter) Status() int {
	return w.finalStatus
}

func (w *WildcardRouterWriter) reset() {
	w.skipCheckStatus = false
	w.finalStatus = 0
	w.tmpStatus = 0
}

func (w WildcardRouterWriter) matchedStatus() bool {
	return w.Status() != http.StatusNotFound && w.Status() != 0
}

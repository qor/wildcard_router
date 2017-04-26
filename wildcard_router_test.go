package wildcard_router_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/qor/wildcard_router"
)

var (
	mux    = http.NewServeMux()
	Server = httptest.NewServer(mux)
)

type ModuleBeforeA struct {
}

func (a ModuleBeforeA) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/module_a0" {
		_, err := w.Write([]byte("Module Before A handled"))
		if err != nil {
			panic("ModuleBeforeA A can't handle")
		}
	} else {
		http.NotFound(w, req)
	}
}

type ModuleA struct {
}

func (a ModuleA) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/module_a0" || req.URL.Path == "/module_a" || req.URL.Path == "/module_ab" {
		_, err := w.Write([]byte("Module A handled"))
		if err != nil {
			panic("Module A can't handle")
		}
	} else {
		http.NotFound(w, req)
	}
}

type ModuleB struct {
}

func (b ModuleB) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/module_b" || req.URL.Path == "/module_ab" {
		_, err := w.Write([]byte("Module B handled"))
		if err != nil {
			panic("Module B can't handle")
		}
	} else {
		http.NotFound(w, req)
	}
}

func init() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("Gin Handle HomePage"))
	})

	wildcardRouter := wildcard_router.New()
	wildcardRouter.MountTo("/", mux)
	wildcardRouter.AddHandler(router)
	wildcardRouter.AddHandler(ModuleBeforeA{})
	wildcardRouter.AddHandler(ModuleA{})
	wildcardRouter.AddHandler(ModuleB{})
	wildcardRouter.NotFoundHandler(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("Sorry, this page was gone!"))
	})
}

type WildcardRouterTestCase struct {
	URL              string
	ExpectStatusCode int
	ExpectHasContent string
}

func TestWildcardRouter(t *testing.T) {
	testCases := []WildcardRouterTestCase{
		{URL: "/", ExpectStatusCode: 200, ExpectHasContent: "Gin Handle HomePage"},
		{URL: "/module_a", ExpectStatusCode: 200, ExpectHasContent: "Module A handled"},
		{URL: "/module_b", ExpectStatusCode: 200, ExpectHasContent: "Module B handled"},
		{URL: "/module_x", ExpectStatusCode: 404, ExpectHasContent: "Sorry, this page was gone!"},
		{URL: "/module_a0", ExpectStatusCode: 200, ExpectHasContent: "Module Before A handled"},
	}

	for i, testCase := range testCases {
		var hasError bool
		req, _ := http.Get(Server.URL + testCase.URL)
		content, _ := ioutil.ReadAll(req.Body)
		if req.StatusCode != testCase.ExpectStatusCode {
			t.Errorf(color.RedString(fmt.Sprintf("WildcardRouter #%v: HTML expect status code '%v', but got '%v'", i+1, testCase.ExpectStatusCode, req.StatusCode)))
			hasError = true
		}
		if string(content) != testCase.ExpectHasContent {
			t.Errorf(color.RedString(fmt.Sprintf("WildcardRouter #%v: HTML expect have content '%v', but got '%v'", i+1, testCase.ExpectHasContent, string(content))))
			hasError = true
		}
		if !hasError {
			fmt.Printf(color.GreenString(fmt.Sprintf("WildcardRouter #%v: Success\n", i+1)))
		}
	}
}

package wildcard_router_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fatih/color"
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
	}
}

func init() {
	WildcardRouter := wildcard_router.New(mux)
	WildcardRouter.AddHandler(ModuleBeforeA{})
	WildcardRouter.AddHandler(ModuleA{})
	WildcardRouter.AddHandler(ModuleB{})
}

type WildcardRouterTestCase struct {
	URL              string
	ExpectHasContent string
}

func TestWildcardRouter(t *testing.T) {
	testCases := []WildcardRouterTestCase{
		{URL: "/module_a", ExpectHasContent: "Module A handled"},
		{URL: "/module_b", ExpectHasContent: "Module B handled"},
		{URL: "/module_x", ExpectHasContent: "404 page not found"},
		{URL: "/module_a0", ExpectHasContent: "Module Before A handled"},
	}

	for i, testCase := range testCases {
		var hasError bool
		req, _ := http.Get(Server.URL + testCase.URL)
		content, _ := ioutil.ReadAll(req.Body)
		if string(content) == testCase.ExpectHasContent {
			t.Errorf(color.RedString(fmt.Sprintf("WildcardRouter #%v: HTML expect have content '%v', but got '%v'", i+1, testCase.ExpectHasContent, string(content))))
			hasError = true
		}
		if !hasError {
			fmt.Printf(color.GreenString(fmt.Sprintf("WildcardRouter #%v: Success\n", i+1)))
		}
	}
}

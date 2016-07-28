package wildcard_router_test

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/qor/wildcard_router"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mux    = http.NewServeMux()
	Server = httptest.NewServer(mux)
)

type ModuleBeforeA struct {
	wildcard_router.WildcardInterface
}

func (a ModuleBeforeA) Handle(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/module_a0" {
		w.Write([]byte("Module Before A handled"))
		return true
	}
	return false
}

type ModuleA struct {
	wildcard_router.WildcardInterface
}

func (a ModuleA) Handle(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/module_a0" {
		w.Write([]byte("Module A handled"))
		return true
	}
	if req.URL.Path == "/module_a" {
		w.Write([]byte("Module A handled"))
		return true
	}
	if req.URL.Path == "/module_ab" {
		w.Write([]byte("Module A handled"))
		return true
	}
	return false
}

type ModuleB struct {
	wildcard_router.WildcardInterface
}

func (b ModuleB) Handle(w http.ResponseWriter, req *http.Request) bool {
	if req.URL.Path == "/module_b" {
		w.Write([]byte("Module B handled"))
		return true
	}
	if req.URL.Path == "/module_ab" {
		w.Write([]byte("Module B handled"))
		return true
	}
	return false
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
		if !strings.Contains(string(content), testCase.ExpectHasContent) {
			t.Errorf(color.RedString(fmt.Sprintf("WildcardRouter #%v: HTML expect have content '%v', but got '%v'", i+1, testCase.ExpectHasContent, string(content))))
			hasError = true
		}
		if !hasError {
			t.Errorf(color.GreenString(fmt.Sprintf("WildcardRouter #%v: Success", i+1)))
		}
	}
}

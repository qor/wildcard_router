# WildcardRouter

WildcardRouter is component that used to handle wildcard routers.

You could choose wildcard_router to satisfy below scenario:

* You have a Module A that could
  - Store a record with URL and Content
  - Have handle and will return content if the URL match one of the records' URL
* You have a Module B have same behaviour as Module A

Using wildcard_router will help you choose correct handler.

[![GoDoc](https://godoc.org/github.com/qor/wildcard_router?status.svg)](https://godoc.org/github.com/qor/wildcard_router)

## Usage

```go
import (
	"github.com/qor/wildcard_router"
)

type ModuleA struct {
	wildcard_router.WildcardInterface
}

func (a ModuleA) Handle(w http.ResponseWriter, req *http.Request) bool {
    // Module A has records:
    //   Record1(URL: /page1, Content: aaa)
    //   Record2(URL: /page2, Content: aaa1)
	if all records' URL contains req.URL.Path {
		w.Write([]byte(aaa or aaa1))
		return true
	}
	return false
}

type ModuleB struct {
	wildcard_router.WildcardInterface
}

func (b ModuleB) Handle(w http.ResponseWriter, req *http.Request) bool {
    // Module B has records:
    //   Record1(URL: /p1, Content: bbb)
    //   Record2(URL: /p2, Content: bbb1)
	if all records' URL contains req.URL.Path {
		w.Write([]byte(bbb or bbb1))
		return true
	}
	return false
}

func main() {
	mux := http.NewServeMux()
	WildcardRouter := wildcard_router.New(mux)
	WildcardRouter.AddHandler(ModuleA{})
	WildcardRouter.AddHandler(ModuleB{})

    // Visit /page1 will return "aaa"
    // Visit /page2 will return "aaa1"
    // Visit /p1 will return "bbb"
    // Visit /p2 will return "bbb1"
    // Visit /unknow will return "404 page not found"
}

```

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).

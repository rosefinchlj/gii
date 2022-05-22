package gii

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	equal := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	if !equal {
		t.Error("parsePattern failed")
	}

	equal = reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	if !equal {
		t.Error("parsePattern failed")
	}

	equal = reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !equal {
		t.Error("parsePattern failed")
	}
}

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)

	return r
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/gii")
	if n == nil {
		t.Fatal("should not be nil")
	}

	if n.pattern != "/hello/:name" {
		t.Fatal("should equal to /hello/:name")
	}

	if ps["name"] != "gii" {
		t.Fatal("name should equal to gii")
	}

	n, ps = r.getRoute("GET", "/assets/css/test.css")
	if n == nil {
		t.Fatal("should not be nil")
	}
}

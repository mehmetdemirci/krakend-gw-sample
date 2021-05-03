package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
)

type registerer string

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer("router-plugin")

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].([]interface{})
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name[0] != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("before router handler %s",html.EscapeString(req.URL.Path))		
		handler.ServeHTTP(w, req)
		fmt.Println("after router-plugin called")
	}), nil
}

func init() {
	fmt.Println("router-plugin handler loaded!!!")
}

func main() {}
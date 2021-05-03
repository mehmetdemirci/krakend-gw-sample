package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

type registerer string

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = registerer("proxy-plugin")

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r registerer) registerClients(ctx context.Context, extra map[string]interface{}) (h http.Handler,e error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("before proxy handler %s",html.EscapeString(req.URL.Path))
		// fmt.Println(req.URL.Query().Get("appid"))

		makeOriginalRequest(w,req)

		fmt.Println("after proxy-plugin called")
	}), nil
}

func init() {
	fmt.Println("proxy-plugin client plugin loaded!!!")
}

func makeOriginalRequest(w http.ResponseWriter, req *http.Request) {
	    client := &http.Client{}
		// Send an HTTP request and returns an HTTP response object.
		 resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
	
		// headers
		for name, values := range resp.Header {
		    w.Header()[name] = values
		}
		
		// status (must come after setting headers and before copying body)
		w.WriteHeader(resp.StatusCode)
		
		body, err := ioutil.ReadAll(resp.Body)		
		w.Write(body)

		fmt.Println("request completed")
}

func main() {}
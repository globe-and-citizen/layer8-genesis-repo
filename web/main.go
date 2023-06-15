package main

import (
	"syscall/js"

	"github.com/globe-and-citizen/layer8-genesis-repo/api"
	"github.com/globe-and-citizen/layer8-genesis-repo/web/protocols/http1"
)

var (
	ServerHost string
	RESTPort   string
	GRPCPort   string
)

func fetch(this js.Value, args []js.Value) interface{} {
	url := args[0].String()
	options := js.ValueOf(map[string]interface{}{
		"method":  "GET",
		"headers": js.ValueOf(map[string]interface{}{}),
	})
	if len(args) > 1 {
		options = args[1]
	}

	method := options.Get("method").String()
	if method == "" {
		method = "GET"
	}
	headers := options.Get("headers")
	body := options.Get("body").String()
	// setting the body to an empty string if it's undefined
	if body == "<undefined>" {
		body = ""
	}

	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// to avoid deadlock with the main thread, we need to run this in a goroutine
		go func() {
			// add headers
			headersMap := make(map[string]string)
			js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				headersMap[args[0].String()] = args[1].String()
				return nil
			}))

			// make the request using the http1 client
			req := api.NewRequest(method, url, headersMap, []byte(body))
			res := http1.NewClient("http", ServerHost, RESTPort).Do(req)
			if res.Status < 300 {
				args[0].Invoke(js.Global().Get("Response").New(string(res.Body), js.ValueOf(map[string]interface{}{
					"status":     res.Status,
					"statusText": res.StatusText,
				})))
			} else {
				args[1].Invoke(js.Global().Get("Error").New(res.StatusText))
			}
		}()
		return nil
	}), js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// for rejection, we just return the error message
		return args[0].String()
	}))

	return promise
}

func main() {
	// keep the thread alive
	close := make(chan struct{}, 0)

	// expose the fetch function to the global scope
	js.Global().Set("layer8", js.ValueOf(map[string]interface{}{
		"fetch": js.FuncOf(fetch),
	}))

	// wait indefinitely
	<-close
}

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"syscall/js"
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
	// setting the body to empty string if it's undefined
	if body == "<undefined>" {
		body = ""
	}

	promise := js.Global().Get("Promise").New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// to avoid deadlock with the main thread, we need to run this in a goroutine
		go func() {
			var (
				req *http.Request
				err error
			)
			if body != "" {
				req, err = http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
			} else {
				req, err = http.NewRequest(method, url, nil)
			}
			if err != nil {
				args[1].Invoke(err.Error())
				return
			}
			// add headers
			js.Global().Get("Object").Call("keys", headers).Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				req.Header.Add(args[1].String(), args[0].String())
				return nil
			}))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				args[1].Invoke(err.Error())
				return
			}
			defer resp.Body.Close()
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				args[1].Invoke(err.Error())
				return
			}
			
			// return the response
			args[0].Invoke(js.Global().Get("Response").New(string(respBody), js.ValueOf(map[string]interface{}{
				"status": resp.StatusCode,
				"statusText": resp.Status,
			})))
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

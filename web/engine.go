package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

type Engine struct {
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.PATH=%q\n", req.URL)
		break
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
		break
	default:
		fmt.Fprintf(w, "404 NOT FOUND:%q\n", req.URL.Path)
	}
}

func engineStarter() {
	//http.HandleFunc("/", indexHandler)
	//http.HandleFunc("/hello", helloHandler)
	//log.Fatal(http.ListenAndServe(":9999", nil))

	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	_, err := fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	if err != nil {
		return
	}
}
func helloHandler(w http.ResponseWriter, req *http.Request) {

	var headers []string
	for k, v := range req.Header {
		headers = append(headers, "Header["+k+"] = "+v[0]+"\n")
	}

	sort.Strings(headers)

	for _, v := range headers {
		_, err := fmt.Fprintf(w, v)
		if err != nil {
			return
		}
	}
	_, err := fmt.Fprintf(w, "hello")
	if err != nil {
		return
	}
}

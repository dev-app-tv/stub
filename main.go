package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	// "github.com/kr/pretty"
)

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := httputil.DumpRequest(r, true)
	fmt.Printf("req=%#v\nerr=%#v\n\n", string(req), err)
	// fmt.Printf("req=%# v\nerr=%# v\n\n", pretty.Formatter(req), pretty.Formatter(err))

	saveDataFile(req, "stub.txt")
	fmt.Printf("write done...\n")

	httputil.NewSingleHostReverseProxy(target)
}

func hijack(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Don't forget to close the connection:
	defer conn.Close()
	bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
	bufrw.Flush()
	s, err := bufrw.ReadString('\n')
	if err != nil {
		log.Printf("error reading string: %v", err)
		return
	}
	fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
	bufrw.Flush()
}

func main() {
	fmt.Println("start listen on :8080")
	http.HandleFunc("/", handler)
	// http.HandleFunc("/", hijack)
	http.ListenAndServe(":8080", nil)
}

var data map[string]Stub

type Stub struct {
	Request  *http.Request
	Response *http.Response
}

func saveDataFile(jsonData []byte, filename string) error {
	err := ioutil.WriteFile(filename, jsonData, 0666)
	return err
}

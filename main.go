package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

func handle(w http.ResponseWriter, r *http.Request) {
	h := strings.TrimSuffix(r.Host, ".academy.dev:8000")

	switch h {
	case "start":
		start(w, r)
	case "bad":
		bad(w, r)
	case "good":
		good(w, r)
	case "wildcard":
		wildcard(w, r)
	case "redirect":
		parts := strings.Split(r.RequestURI, "/")

		rhost := strings.TrimLeft(r.RequestURI, "/")

		if len(parts) > 2 {
			rhost = parts[1]
		}

		rlocation := "http://" + rhost + ".academy.dev:8000"

		if len(parts) > 2 {
			rlocation += "/" + parts[2]
		}

		redirect(w, r, rlocation)
	}
}

func start(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintf(w, "%s", dump)
}

func good(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://start.academy.dev:8000")
	w.Header().Set("Access-Control-Allow-Headers", "X-Custom-Header")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if strings.Contains(r.RequestURI, "cookie") {
		w.Header().Set("Set-Cookie", "foo=Server set this cookie")
	}

	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintf(w, "%s", dump)
}

func wildcard(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "X-Custom-Header")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if strings.Contains(r.RequestURI, "cookie") {
		w.Header().Set("Set-Cookie", "foo=Server set this cookie")
	}

	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintf(w, "%s", dump)
}

func redirect(w http.ResponseWriter, r *http.Request, rh string) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "OPTIONS" {
		http.Redirect(w, r, rh, 302)
	} else {
		dump, _ := httputil.DumpRequest(r, true)
		fmt.Fprintf(w, "%s", dump)
	}
}

func bad(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Fprintf(w, "%q", dump)
}

func main() {
	http.HandleFunc("/", handle)

	fs := http.FileServer(http.Dir("."))
	http.Handle("/f/", http.StripPrefix("/f/", fs))

	http.ListenAndServe(":8000", nil)
}

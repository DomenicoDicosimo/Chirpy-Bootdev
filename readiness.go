package main

import "net/http"

func myhandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	body := ([]byte)("OK")
	w.Write(body)
}

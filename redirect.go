package main

import (
	"log"
	"net/http"
)

type redirectHandler struct {}
func (rh *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "https"
	//log.Println(r.URL.String())
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
}


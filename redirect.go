package main

import (
	"net/http"
)

type redirectHandler struct {}
func (rh *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "https"
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
}


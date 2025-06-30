package handlers

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateRouters(router *mux.Router) {

	router.HandleFunc("/hello", printHello())

}

func printHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello, world!")
	}
}

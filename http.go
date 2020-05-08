package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func initHandler() {
	http.HandleFunc("/ping", handlePing)
	// TODO: Adding more APIs

	router := httprouter.New()

	router.GET("/get_prize", handleGetPrize)

	fmt.Println("serving in 8181 ")
	http.ListenAndServe(":8181", nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func handleGetPrize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Write([]byte("pong"))
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/json"
)

func initHandler() {
	http.HandleFunc("/ping", handlePing)
	http.HandleFunc("/get_prize", handleGetPrize)
	http.HandleFunc("/pool", handleGetPrize)
	// TODO: Adding more APIs

	fmt.Println("serving in 8181 ")
	http.ListenAndServe(":8181", nil)
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func handleGetPrize(w http.ResponseWriter, r *http.Request) {
	// userID := r.FormValue("user_id")

	// get redis by user_id
	ovo, err := redisGetRand()
	if err != nil {
		fmt.Println(err)
	}

	resp, err := json.Marshal(ovo)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(resp)
}

func handlePostPrize(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}

	ovos := []KuponOVO{}

	err = json.Unmarshal([]byte(body), ovos)

	for _, ovo := range ovos {
		redisAddKey(ovo)
	}
}

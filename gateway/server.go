package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type LeaderAddress struct {
	M    sync.Mutex
	Addr string
}

var leaderAddr = LeaderAddress{sync.Mutex{}, ""}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/read/{key}", GetValue).Methods("GET")
	r.HandleFunc("/delete/{key}", DeleteValue).Methods("DELETE")
	r.HandleFunc("/create/{key}/{value}", PostValue).Methods("POST")
	r.HandleFunc("/update/{key}/{value}", UpdateValue).Methods("PUT")

	r.HandleFunc("/updateLeaderAddress", UpdateLeaderAddr).Methods("POST")

	http.ListenAndServe(":8070", r)

}

type AddrInform struct {
	Addr string `json:"addr"`
	Port string `json:"port"`
}

func UpdateLeaderAddr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var addrJSON AddrInform
	err := json.NewDecoder(r.Body).Decode(&addrJSON)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}
	addr := addrJSON.Addr
	leaderAddr.M.Lock()
	leaderAddr.Addr = addr + ":" + addrJSON.Port
	leaderAddr.M.Unlock()
	log.Println(addr + " IS THE NEW  UPDATED LEADER!!!")
	fmt.Fprint(w, http.StatusOK)
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	leaderAddr.M.Lock()
	addr := leaderAddr.Addr
	leaderAddr.M.Unlock()
	// http.Redirect(w, r, "http://"+addr+r.URL.String(), 303)
	req, err := http.NewRequest(http.MethodGet, "http://"+addr+r.URL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprint(w, string(resBody))

}
func PostValue(w http.ResponseWriter, r *http.Request) {
	leaderAddr.M.Lock()
	addr := leaderAddr.Addr
	leaderAddr.M.Unlock()
	// http.Redirect(w, r, "http://"+addr+r.URL.String(), 303)
	req, err := http.NewRequest(http.MethodPost, "http://"+addr+r.URL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprint(w, string(resBody))

}

func UpdateValue(w http.ResponseWriter, r *http.Request) {
	leaderAddr.M.Lock()
	addr := leaderAddr.Addr
	leaderAddr.M.Unlock()
	// http.Redirect(w, r, "http://"+addr+r.URL.String(), 303)
	req, err := http.NewRequest(http.MethodPut, "http://"+addr+r.URL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprint(w, string(resBody))
}
func DeleteValue(w http.ResponseWriter, r *http.Request) {
	leaderAddr.M.Lock()
	addr := leaderAddr.Addr
	leaderAddr.M.Unlock()
	// http.Redirect(w, r, "http://"+addr+r.URL.String(), 303)
	req, err := http.NewRequest(http.MethodDelete, "http://"+addr+r.URL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	res, err := http.DefaultClient.Do(req)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprint(w, string(resBody))
}

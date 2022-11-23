package datastore

// !Mainly used for Partition Leader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

var serversMap map[int]string

type RoundBoutCounter struct {
	m sync.Mutex
	c int
}

var rbc = RoundBoutCounter{sync.Mutex{}, 1}

func GetRouter(m map[int]string) *mux.Router {
	serversMap = m
	r := mux.NewRouter()
	r.HandleFunc("/read/{key}", GetValue).Methods("GET")
	r.HandleFunc("/delete/{key}", DeleteValue).Methods("DELETE")
	r.HandleFunc("/create/{key}/{value}", PostValue).Methods("POST")
	r.HandleFunc("/update/{key}/{value}", UpdateValue).Methods("PUT")
	return r
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	log.Println("AUFF GET")
	vars := mux.Vars(r)
	key := vars["key"]

	var resp string
	for i := 1; i <= len(serversMap); i++ {
		res, err := DialTCPServer(serversMap[i], key, "NONE", "GET")
		if err != nil {
			continue
		}
		if res != "NOTFOUND" {
			resp += res
			break
		}
	}
	fmt.Fprint(w, resp)

}

func DeleteValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var resp string
	for i := 1; i <= len(serversMap); i++ {
		res, err := DialTCPServer(serversMap[i], key, "NONE", "DELETE")
		if err != nil {
			continue
		}
		resp += res
	}

	fmt.Fprint(w, resp)
}

func UpdateValue(w http.ResponseWriter, r *http.Request) {
	log.Println("AUFFF GET")
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	var resp string
	for i := 1; i <= len(serversMap); i++ {
		res, err := DialTCPServer(serversMap[i], key, val, "PUT")
		if err != nil {
			continue
		}
		resp = res
	}
	fmt.Fprint(w, resp)

}
func PostValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	rbc.m.Lock()
	if rbc.c%(len(serversMap)+1) == 0 {
		rbc.c = 1
	}
	temp := rbc.c
	rbc.c += 1
	rbc.m.Unlock()
	var resp string
	for i := 0; i < int(len(serversMap)/2+1); i++ {
		res, err := DialTCPServer(serversMap[temp], key, val, "POST")
		if err != nil {
			continue
		}
		resp += res
		temp += 1
		if temp%(len(serversMap)+1) == 0 {
			temp = 1
		}
	}
	fmt.Fprint(w, resp)
}

func DialTCPServer(tcp_addr string, key, val, cmd string) (string, error) {
	conn, err := net.Dial("tcp", tcp_addr)
	if err != nil {
		return "", err
	}
	msgStruct := TCPMsg{Cmd: cmd, Key: key, Val: val}
	msg, _ := json.Marshal(msgStruct)
	fmt.Fprint(conn, string(msg))
	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	return message, nil

}

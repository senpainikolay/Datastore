package datastore

// !Mainly used for Partition Leader

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	r.HandleFunc("/create/{key}/{value}", PostValue).Methods("POST")
	return r
}

func GetValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	resp := DialTCPServer(serversMap[1], key, "NONE", "GET")
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
		resp += DialTCPServer(serversMap[temp], key, val, "POST")
		temp += 1
		if temp%(len(serversMap)+1) == 0 {
			temp = 1
		}
	}
	fmt.Fprint(w, resp)
}

func DialTCPServer(tcp_addr string, key, val, cmd string) string {
	conn, err := net.Dial("tcp", tcp_addr)
	if err != nil {
		fmt.Println("error:", err)
	}
	msgStruct := TCPMsg{Cmd: cmd, Key: key, Val: val}
	msg, _ := json.Marshal(msgStruct)
	fmt.Fprint(conn, string(msg))
	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	return message

}

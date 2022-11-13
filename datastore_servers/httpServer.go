package datastore

// !Mainly used for Partition Leader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

var serversMap map[int]string

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
	resp := DialTCPServerOnGet(serversMap[1], key)
	fmt.Fprint(w, resp)

}
func PostValue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	resp := DialTCPServerOnPost(serversMap[1], key, val)
	fmt.Fprint(w, resp)

}

func DialTCPServerOnGet(tcp_addr string, key string) string {
	conn, err := net.Dial("tcp", tcp_addr)
	if err != nil {
		fmt.Println("error:", err)
	}
	msgStruct := TCPMsg{Cmd: "GET", Key: key}
	msg, _ := json.Marshal(msgStruct)
	fmt.Fprint(conn, string(msg))

	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	return message

}
func DialTCPServerOnPost(tcp_addr string, key, val string) string {
	conn, err := net.Dial("tcp", tcp_addr)
	if err != nil {
		fmt.Println("error:", err)
	}
	msgStruct := TCPMsg{Cmd: "POST", Key: key, Val: val}
	msg, _ := json.Marshal(msgStruct)
	fmt.Fprint(conn, string(msg))

	message, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
	return message

}

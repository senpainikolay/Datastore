package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var Map sync.Map

func RunTCPServer(tcp_addr string) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", tcp_addr)
	if err != nil {
		log.Panicln("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Printf("Listening on %v \n", tcp_addr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, tcp_addr)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn, tcp_addr string) {
	// we create a decoder that reads directly from the socket
	d := json.NewDecoder(conn)
	var msg TCPMsg
	err := d.Decode(&msg)
	if err != nil {
		log.Fatal("couldnt decode the TCP server json format")
	}
	var resp string
	switch msg.Cmd {
	case "POST":
		Map.Store(msg.Key, msg.Val)
		resp = fmt.Sprintf("Added/Posted at %v", tcp_addr)

	case "GET":
		storeVal, ok := Map.Load(msg.Key)
		if ok {
			resp = fmt.Sprint(storeVal)
		} else {
			resp = "NOTFOUND"
		}
	case "DELETE":
		_, ok := Map.Load(msg.Key)
		if ok {
			Map.Delete(msg.Key)
			resp = fmt.Sprintf("Deleted at %v", tcp_addr)
		}
	case "PUT":
		_, ok := Map.Load(msg.Key)
		if ok {
			Map.Store(msg.Key, msg.Val)
		}

	}
	conn.Write([]byte(resp))
	conn.Close()
}

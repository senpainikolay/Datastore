package datastore

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func RunTCPServer(tcp_addr string) {
	// Listen for incoming connections.
	l, err := net.Listen("tcp", tcp_addr)
	if err != nil {
		log.Panicln("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	//log.Printf("Listening on %v \n", tcp_addr)
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
		MapM.M.Lock()
		MapM.Map[msg.Key] = msg.Val
		MapM.M.Unlock()
		go func() {
			MapLenCounter.M.Lock()
			MapLenCounter.C += 1
			MapLenCounter.M.Unlock()
		}()

		resp = fmt.Sprintf("Added/Posted at %v", tcp_addr)

	case "GET":
		MapM.M.Lock()
		storeVal, ok := MapM.Map[msg.Key]
		MapM.M.Unlock()
		if ok {
			resp = fmt.Sprint(storeVal)
		} else {
			resp = "NOTFOUND"
		}
	case "DELETE":
		MapM.M.Lock()
		_, ok := MapM.Map[msg.Key]
		MapM.M.Unlock()
		if ok {
			MapM.M.Lock()
			delete(MapM.Map, msg.Key)
			MapM.M.Unlock()
			go func() {
				MapLenCounter.M.Lock()
				MapLenCounter.C -= 1
				MapLenCounter.M.Unlock()
			}()
			resp = fmt.Sprintf("Deleted at %v", tcp_addr)
		}
	case "PUT":
		MapM.M.Lock()
		_, ok := MapM.Map[msg.Key]
		MapM.M.Unlock()

		if ok {
			MapM.M.Lock()
			MapM.Map[msg.Key] = msg.Val
			MapM.M.Unlock()
			resp = fmt.Sprintf("Updated at %v", tcp_addr)
		}

	}
	conn.Write([]byte(resp))
	conn.Close()
}

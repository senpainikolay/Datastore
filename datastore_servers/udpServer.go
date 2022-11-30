package datastore

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
)

type MapLenCounterType struct {
	M sync.Mutex
	C int
}

var MapLenCounter = MapLenCounterType{sync.Mutex{}, 0}

func sendUDPResponse(conn *net.UDPConn, addr *net.UDPAddr, resp string) {
	_, err := conn.WriteToUDP([]byte(resp), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func RunUDPServer(port string, dsAddr string) {
	portInt, _ := strconv.Atoi(port)
	p := make([]byte, 2048)
	addr := net.UDPAddr{
		Port: portInt,
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			return
		}

		keyMap := ParseByteArr(p)
		found := 0
		MapM.M.Lock()
		_, ok := MapM.Map[keyMap]
		MapM.M.Unlock()
		log.Println(ok)
		if ok {
			found = 1
		}

		MapLenCounter.M.Lock()
		tempLen := MapLenCounter.C
		MapLenCounter.M.Unlock()
		respAddr := fmt.Sprintf("%v:%v", dsAddr, port)

		respString := fmt.Sprintf("%v %v %v", found, tempLen, respAddr)

		go sendUDPResponse(ser, remoteaddr, respString)
	}
}

func ParseByteArr(p []byte) string {
	var kek []byte
	for i := range p {
		if p[i] == 0 {
			break
		}
		kek = append(kek, p[i])
	}
	return string(kek)
}

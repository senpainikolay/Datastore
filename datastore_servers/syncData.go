package datastore

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func SyncTemporarily(clusterServers map[int]string) {

	for {

		time.Sleep(time.Duration((rand.Intn(5))+10) * time.Second)

		MapM.M.Lock()
		tempMap := MapM.Map
		MapM.M.Unlock()
		for key, _ := range tempMap {
			// List of Responses from UDP : Servers that have the actual Key
			var resHave [][]string
			// List of Responses from UDP : Servers that DO NOT have the actual Key
			var resLack [][]string
			for i := 2; i <= len(clusterServers); i++ {
				resString := DialUDP(clusterServers[i], key)
				// if server not dead
				if resString != "" {
					resSlice := strings.Split(resString, " ")
					// [0] : Bool: have key/do not have key
					// [1] : Lenght of the Database Map.
					// [2] : The addreass
					if resSlice[0] == "1" {
						resHave = append(resHave, resSlice)
					} else {
						resLack = append(resLack, resSlice)
					}

				}
			}

			log.Println(len(resHave))
			log.Println(resLack)
			// filter the servers
			for {
				// Other servers that have that key +  +1:actual server.
				if len(resHave)+1 > len(resLack) || len(resLack) == 1 {
					break
				}
				var tempMinMapLen = 9999
				var idx = 0
				// Select server with smallest Map and adding/completing Data loss on server restart. ( in a while )
				for i, item := range resLack {
					intLen, _ := strconv.Atoi(item[1])
					if intLen < tempMinMapLen {
						tempMinMapLen = intLen
						idx = i
					}
				}
				syncRes, _ := DialTCPServer(resLack[idx][2], key, tempMap[key], "POST")
				log.Printf("Syncronizing %v  with key %v  and the response: %s \n", resLack[idx][2], key, syncRes)
				resLack = removeElemByIndex(resLack, idx)
			}
		}
	}
}

func DialUDP(addr string, key string) string {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return ""
	}
	fmt.Fprint(conn, key)
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		conn.Close()
		return ParseByteArr(p)
	} else {
		fmt.Printf("Some error %v\n", err)
		conn.Close()
		return ""
	}

}

func removeElemByIndex(slice [][]string, s int) [][]string {
	return append(slice[:s], slice[s+1:]...)
}

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
			var res [][]string
			var deadServers = 0
			for i := 1; i <= len(clusterServers); i++ {
				resString := DialUDP(clusterServers[i], key)
				if resString != "" {
					resSlice := strings.Split(resString, " ")
					res = append(res, resSlice)

				} else {
					deadServers += 1

				}
			}
			// 2+1 -deadServers :1have       1:0no

			log.Println(res)

			// filter the servers
			for {
				if deadServers == len(clusterServers)-1 || len(res) <= int((len(clusterServers)-deadServers-1)/2) {
					break
				}
				var tempMinMapLen = 0
				var idx = 0
				for i, item := range res {
					intLen, _ := strconv.Atoi(item[1])
					if intLen > tempMinMapLen {
						tempMinMapLen = intLen
						idx = i
					}
				}
				res = removeElemByIndex(res, idx)
			}
			// Updating those servers with data.
			for _, item := range res {
				val := tempMap[key]
				syncRes, _ := DialTCPServer(item[2], key, val, "POST")
				log.Printf("Synconized  %v  with key %v  and the response: %s \n", item[2], key, syncRes)
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

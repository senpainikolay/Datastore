package main

import (
	"log"
	"net/http"

	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

func main() {

	conf := GetConf()
	if conf.LeaderBool {
		r := datastore.GetRouter(conf.ServerMap)
		go datastore.RunTCPServer(conf.TcpAddr)
		log.Printf("THE LEADER IS: %v : %v ", conf.HttpAddr, conf.HttpPort)
		http.ListenAndServe(":"+conf.HttpPort, r)
	} else {
		datastore.RunTCPServer(conf.TcpAddr)
	}

}

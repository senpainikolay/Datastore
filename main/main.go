package main

import (
	"log"
	"net/http"

	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

func main() {

	conf := GetConf()
	r := datastore.GetRouter(conf.ServerMap)
	log.Println(conf)
	go datastore.RunTCPServer(conf.TcpAddr)
	http.ListenAndServe(":"+conf.HttpPort, r)

}

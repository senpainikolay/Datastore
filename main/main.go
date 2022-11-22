package main

import (
	"net/http"

	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

func main() {
	var conf = GetConf()
	r := datastore.GetRouter(conf.ServerMap)
	r = AttachLeaderFeatureToHTTP(r)
	go datastore.RunTCPServer(conf.TcpAddr)
	go CompeteForLeader()
	http.ListenAndServe(":"+conf.HttpPort, r)

}

package main

import (
	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

func main() {
	go datastore.RunTCPServer(conf.TcpAddr)
	AttachLeaderFeatureToHTTP()

}

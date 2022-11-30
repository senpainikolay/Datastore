package main

import (
	"fmt"

	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

func main() {
	go datastore.RunTCPServer(fmt.Sprintf("%v:%v", conf.Addr, conf.TcpPort))
	go datastore.RunUDPServer(conf.TcpPort, conf.Addr)
	go datastore.SyncTemporarily(conf.ServerMap)
	AttachLeaderFeatureToHTTP()

}

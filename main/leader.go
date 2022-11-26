package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	datastore "github.com/senpainikolay/Datastore/datastore_servers"
)

type Leader struct {
	M sync.Mutex
	B bool
}

var leader = Leader{sync.Mutex{}, false}
var conf = GetConf()
var SleepControllerChan = make(chan int, 5)
var LeaderConfirmationChan = make(chan int, 1)
var LeaderChan = make(chan int, 3)

func AttachLeaderFeatureToHTTP() {
	r := datastore.GetRouter(conf.ServerMap)
	r.HandleFunc("/leaderInfo", GetLeaderInfo).Methods("GET")
	r.HandleFunc("/dominate", DominatedFn).Methods("GET")
	go CompeteForLeader()
	http.ListenAndServe(":"+conf.HttpPort, r)
}

func GetLeaderInfo(w http.ResponseWriter, r *http.Request) {
	go func() { SleepControllerChan <- 100 }()
	leader.M.Lock()
	isLeader := leader.B
	leader.M.Unlock()
	var resp int8
	if isLeader {
		resp = 1
	}

	fmt.Fprint(w, resp)

}

func DominatedFn(w http.ResponseWriter, r *http.Request) {
	go func() { SleepControllerChan <- 200 }()
	fmt.Fprint(w, "YOU ARE THE LEADER!")

}

func CompeteForLeader() {

	time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)

	for {

		select {

		case timeUnit := <-SleepControllerChan:
			time.Sleep(time.Duration(timeUnit) * time.Millisecond)

		case leaderSleepTime := <-LeaderChan:
			go func() { LeaderChan <- 200 }()
			go DominateOnTimeServers()
			time.Sleep(time.Duration(leaderSleepTime) * time.Millisecond)

		case <-LeaderConfirmationChan:
			if CheckLeaderStats() != 0 {
				leader.M.Lock()
				leader.B = false
				leader.M.Unlock()
			} else {
				go func() { LeaderChan <- 150 }()
				go InformGatewayServer()
				DominateOnTimeServers()
			}

		default:
			if CheckLeaderStats() == 0 {

				leader.M.Lock()
				leader.B = true
				leader.M.Unlock()
				go func() { LeaderConfirmationChan <- -1 }()
			}

		}
	}

}

func CheckLeaderStats() int {
	var wg sync.WaitGroup
	wg.Add(len(conf.HttpServerMap))
	c := 0
	for _, v := range conf.HttpServerMap {
		temp := v
		go func() {
			c += CheckLeadederInfo(temp)
			wg.Done()
		}()
	}
	wg.Wait()
	return c
}

func CheckLeadederInfo(addr string) int {
	resp, err := http.Get("http://" + addr + "/leaderInfo")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if string(body) == "1" {
		return 1
	}
	return 0

}

func DominateOnTimeServers() {
	for _, v := range conf.HttpServerMap {
		temp := v
		go func() {
			AddTimeOnServerAddress(temp)
		}()
	}
}

func AddTimeOnServerAddress(addr string) {
	_, err := http.Get("http://" + addr + "/dominate")
	if err != nil {
		return
	}
}

func InformGatewayServer() {
	postBody, _ := json.Marshal(map[string]string{
		"addr": conf.HttpAddr,
		"port": conf.HttpPort,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://gateway:8070/updateLeaderAddress", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	resp.Body.Close()

}

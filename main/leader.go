package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
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

func AttachLeaderFeatureToHTTP(r *mux.Router) *mux.Router {
	r.HandleFunc("/leaderInfo", GetLeaderInfo).Methods("GET")
	r.HandleFunc("/dominate", DominatedFn).Methods("GET")
	r.HandleFunc("/a", B).Methods("GET")
	return r
}
func B(w http.ResponseWriter, r *http.Request) {
	leader.M.Lock()
	resp := leader.B
	leader.M.Unlock()

	fmt.Fprint(w, resp)

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
				DominateOnTimeServers()
				log.Println(conf.HttpAddr + " is the LEADER!!!")
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

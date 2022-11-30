package datastore

import "sync"

type TCPMsg struct {
	Cmd string `json:"cmd"`
	Val string `json:"val"`
	Key string `json:"key"`
}

type MapMutex struct {
	M   sync.Mutex
	Map map[string]string
}

var MapM = MapMutex{sync.Mutex{}, make(map[string]string, 0)}

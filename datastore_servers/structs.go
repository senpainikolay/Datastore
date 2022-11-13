package datastore

type TCPMsg struct {
	Cmd string `json:"cmd"`
	Val string `json:"val"`
	Key string `json:"key"`
}

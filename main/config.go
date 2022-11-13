package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Conf struct {
	HttpPort   string         `json:"http_port"`
	TcpAddr    string         `json:"tcp_address"`
	LeaderBool bool           `json:"leader_bool"`
	ServerMap  map[int]string `json:"tcp_cluster_servers"`
}

func GetConf() *Conf {
	jsonFile, err := os.Open("config/config.json")
	if err != nil {
		log.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Conf
	json.Unmarshal(byteValue, &conf)
	return &conf
}

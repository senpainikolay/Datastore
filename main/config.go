package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Conf struct {
	HttpPort      string         `json:"http_port"`
	Addr          string         `json:"addr"`
	TcpPort       string         `json:"tcp_port"`
	ServerMap     map[int]string `json:"cluster_servers"`
	HttpServerMap map[int]string `json:"http_servers"`
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

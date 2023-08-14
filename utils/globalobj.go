package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type GlobalObj struct {
	Host             string
	TcpProt          int
	Name             string
	Version          string
	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxConn          int
	MaxPackageSize   uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Name:             "App",
		Version:          "V0.4",
		TcpProt:          8888,
		Host:             "0.0.0.0",
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 512,
		MaxConn:          1024,
		MaxPackageSize:   1024,
	}
	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/conf.json")
	if err != nil {
		fmt.Println("not loading conf.json")
		return
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func (g *GlobalObj) GetName() string {
	return g.Name
}

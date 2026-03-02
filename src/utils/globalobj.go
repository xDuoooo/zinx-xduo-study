package utils

import (
	"encoding/json"
	"os"
	"zinx-xduo-study/src/ziface"
)

// GlobalObj 全局配置类
type GlobalObj struct {
	TCPServer ziface.IServer //全局Server对象
	Host      string
	TCPPort   int
	Name      string

	Version        string //Zinx 版本号
	MaxConn        int    //允许的最大连接数
	MaxPackageSize uint32 //数据包的最大值
}

//全局对外的Glovalobj

var GlobalObject *GlobalObj

//从配置文件加载

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//json文件解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "0.4",
		TCPPort:        8999,
		MaxConn:        1000,
		MaxPackageSize: 4096,
		Host:           "0.0.0.0",
	}
	//GlobalObject.Reload()
}

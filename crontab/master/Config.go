package master

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriteTimeout int `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int	`json:"etcdDialTimeout"`
	WebRoot string `json:"webroot"`
}

var(
	// 单例
	G_config *Config
)

// 加载配置
func InitConfig(filename string)(err error){
	var(
		content []byte
		config Config
	)
	// 1 把配置文件读进来
	if content,err = ioutil.ReadFile(filename); err != nil{
		return
	}

	// 2 JSON反序列化
	if err = json.Unmarshal(content,&config); err != nil{
		return
	}

	// 3 赋值单例
	G_config = &config

	//fmt.Println(config)

	return
}
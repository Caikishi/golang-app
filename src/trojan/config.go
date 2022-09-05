package trojan

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Port      string           `yaml:"port"`
	SocksPort string           `yaml:"socks-port"`
	Proxies   []map[string]any `yaml:"proxies"`
}

func init() {

}

// 读取Yaml配置文件,
// 并转换成conf对象
func GetConf() *conf {
	//应该是 绝对地址
	yamlFile, err := ioutil.ReadFile("config2.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}

	var c *conf
	err = yaml.Unmarshal(yamlFile, &c)

	if err != nil {
		fmt.Println(err.Error())
	}

	return c
}

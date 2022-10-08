package trojan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type myJson struct {
	Outbounds []streams `json:"outbounds"`
	Routing   any       `json:"routing"`
}

type streams struct {
	Tag         string     `json:"tag"`
	Protocol    string     `json:"protocol"`
	Settings    mySettings `json:"settings"`
	SendThrough string     `json:"sendThrough"`
}

type mySettings struct {
	Servers []MyServers `json:"servers"`
}

type MyServers struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Method   string `json:"method"`
	Password string `json:"password"`
}

func GetJson(d []MyServers) {
	data, err := ioutil.ReadFile("test.json")
	if err != nil {
		return
	}
	v := myJson{}
	fmt.Printf("v: %v\n", v)
	json.Unmarshal(data, &v)

	for _, v := range v.Outbounds {
		if v.Tag == "jp" {
			v.Settings.Servers = nil
			for _, k := range d {
				fmt.Printf("k: %v\n", k)
			}
			fmt.Printf("v: %v\n", v.Settings.Servers)
		}
	}
	j, _ := json.MarshalIndent(v, "", "  ")
	ioutil.WriteFile("test2.json", j, 0777)

}

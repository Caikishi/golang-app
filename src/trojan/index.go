package trojan

import (
	"io/ioutil"
	"net/http"
)

var url = ""

func GetYaml() {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("config2.yaml", data, 0644)
}

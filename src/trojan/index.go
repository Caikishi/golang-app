package trojan

import (
	"io/ioutil"
	"net/http"
)

var url = "https://api.subcloud.xyz/sub?target=clash&url=https%3A%2F%2Fsubscription.kyapi.xyz%2Fmodules%2Fservers%2FV2RaySocks%2Fsubscribe%2Fv2rayng.php%3Fsid%3D108439%26token%3DArNpiWCAEKPs&insert=false&config=https%3A%2F%2Fraw.githubusercontent.com%2FACL4SSR%2FACL4SSR%2Fmaster%2FClash%2Fconfig%2FACL4SSR_Online.ini&exclude=%E5%9B%9E%E5%9B%BD%E8%B7%AF%E7%BA%BF%7C%E4%BB%85%E6%B5%B7%E5%A4%96%E7%94%A8%E6%88%B7%7C%E5%89%A9%E4%BD%99%E6%B5%81%E9%87%8F%7C%E5%88%B0%E6%9C%9F%E6%97%B6%E9%97%B4%7C%E5%80%8D%E7%8E%87&filename=wenyun&emoji=true&list=false&tfo=false&scv=false&fdn=false&sort=false&new_name=true"

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

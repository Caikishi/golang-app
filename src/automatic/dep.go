package automatic

import (
	"bytes"
	"example/src/gee"
	"fmt"
	"os"
	"os/exec"
)

func BuildReact(ctx *gee.Context) {
	fmt.Println("React push 开始验证")
	if token := ctx.Req.Header["X-Gitee-Token"]; token[0] != "feixun@123" {
		return
	}
	fmt.Println("验证通过,开始拉取代码")
	url := "/Users/caikishi/Documents/CODE/my-react-app/"
	pullGit(url)
	yarnBuild(url)
}

func pullGit(url string) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = url
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		fmt.Printf("stderr.String(): %v\n", stderr.String())
	}
}

func yarnBuild(url string, a ...string) {
	cmd := exec.Command("yarn", "build")
	cmd.Dir = url
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		fmt.Printf("stderr.String(): %v\n", stderr.String())
	}
}

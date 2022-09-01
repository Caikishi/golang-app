package automatic

import (
	"example/src/gee"
	"fmt"
	"log"
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
	fmt.Println("start build")
	yarnBuild(url)
}

func pullGit(url string) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = url
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}

}

func yarnBuild(url string, a ...string) {
	cmd := exec.Command("yarn", "build")
	cmd.Dir = url
	fmt.Println()
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
}

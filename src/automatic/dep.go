package automatic

import (
	"bytes"
	"example/src/gee"
	"example/src/geeconfig"
	"fmt"
	"os"
	"os/exec"
)

var config = geeconfig.GetConf()
var flag = make(chan bool)
var javaCount int
var reactCount int

func BuildReact(ctx *gee.Context) {
	defer func() {
		reactCount = 0
	}()
	if !verification(ctx) {
		return
	}
	//TODO
	//简单防止多个同时执行(不能解决并发)
	if reactCount > 0 {
		return
	}
	reactCount++
	fmt.Println("验证通过,开始拉取代码 feixun-bs-web")
	pullGit(config.ReactUrl)
	yarn(config.ReactUrl)
	yarnBuild(config.ReactUrl)
}

func BuildJava(ctx *gee.Context) {
	if !verification(ctx) {
		return
	}
	fmt.Println("验证通过,开始拉取代码 java")
	if javaCount > 0 {
		flag <- true
	}
	pullGit(config.JavaRul)
	fmt.Println("开始打包代码")
	mvnPackage(config.JavaRul)
	fmt.Println("开始打包代码")
	javaRun(config.JavaRul)

}

func verification(ctx *gee.Context) bool {
	fmt.Printf("config: %v\n", config)
	fmt.Println("React push 开始验证密码")
	token := ctx.Req.Header["X-Gitee-Token"]
	if len(token) == 0 || token[0] != config.Password {
		fmt.Println("密码验证不通过")
		return false
	}
	return true
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
func yarn(url string) {
	cmd := exec.Command("yarn")
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
func yarnBuild(url string) {
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

func mvnPackage(url string) {
	cmd := exec.Command("mvn", "package")
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

func javaRun(url string) {
	cmd := exec.Command("java", "-jar", "demo-0.0.1-SNAPSHOT.jar")
	cmd.Dir = url
	var stderr bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	javaCount++
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		fmt.Printf("stderr.String(): %v\n", stderr.String())
	}
	v := <-flag
	if v {
		cmd.Process.Kill()
		javaCount = 0
	}

}

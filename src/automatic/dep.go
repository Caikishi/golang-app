package automatic

import (
	"bytes"
	"example/src/gee"
	"example/src/geeconfig"
	"fmt"
	"log"
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
	// pullGit(config.JavaRul)
	fmt.Println("开始打包代码")
	// mvnPackage(config.JavaRul)
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
	// var stderr bytes.Buffer
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = &stderr
	stdout, err := cmd.StdoutPipe()
	if err = cmd.Start(); err != nil {
		return
	}
	// err := cmd.Start()
	javaCount++
	errorFlag := true
	go func() {
		v := <-flag
		if v {
			//设置自行触发的结束不会进入 error
			errorFlag = false
			cmd.Process.Kill()
			javaCount = 0
		}
	}()
	go func() {
		cmd.Process.Wait()
		if errorFlag {
			if err != nil {
				log.Fatalf("failed to call cmd.Start(): %v", err)
			}
			//TODO
			//后续可能将错误记录直接记录到 redis
			//kill 进程也会触发
			log.Printf("exitcode: %d", cmd.ProcessState.ExitCode())
		} else {
			errorFlag = true
		}
	}()
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		log.Print(string(tmp))
		if err != nil {
			break
		}
	}

}

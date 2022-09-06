package automatic

import (
	"bytes"
	"example/src/gee"
	"example/src/geeconfig"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
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
	pullGit(config.ReactUrl)
	yarn(config.ReactUrl)
	yarnBuild(config.ReactUrl)
}

func BuildJava(ctx *gee.Context) {
	if !verification(ctx) {
		return
	}
	if javaCount > 0 {
		flag <- true
	}
	pullGit(config.JavaRul)
	mvnPackage(config.JavaRul)
	javaRun(config.JavaRul)

}

func verification(ctx *gee.Context) bool {
	fmt.Printf("config: %v\n", config)
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

// 判断文件或文件夹是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
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
	dir := "javaLogs"
	if !isExist(dir) {
		os.Mkdir(dir, 0777)
	}

	cmd := exec.Command("java", "-jar", "feixun-web.jar")
	cmd.Dir = url + "feixun-web/target"
	fmt.Printf("cmd.Dir: %v\n", cmd.Dir)
	stdout, err := cmd.StdoutPipe()
	if err = cmd.Start(); err != nil {
		return
	}
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
			//后续可能将错误记录直接记录到 文件 或 数据库
			log.Printf("exitcode: %d", cmd.ProcessState.ExitCode())
		} else {
			errorFlag = true
		}
	}()
	time := time.Now().Format(dir+"/2006-01-02 15-04-05") + ".log"
	f, err := os.Create(time)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		f.Write([]byte(string(tmp)))
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

}

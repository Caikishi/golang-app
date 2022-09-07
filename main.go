package main

import (
	"example/src/automatic"
	"example/src/gee"
	"fmt"
	"log"
)

type Person interface {
	getName() string
}

type Student struct {
	name string
	age  int
}

// func (stu *Student) getName() string {
// 	return stu.name
// }

type Worker struct {
	name   string
	gender string
}

func (w *Worker) getName() string {
	fmt.Println("我是 worker 我处理一些 worker 在 getName 的逻辑")
	return w.name
}

func ex(p Person) {
	fmt.Println("我是 ex 我处理一些 ex 的逻辑 然后调用了接口 getName")
	fmt.Println(p.getName())

}
func ex2(s *Student) {
	fmt.Println(s.name)
}

type id int64

func Sum[T id](args ...T) T {
	var sum T
	for i := 0; i < len(args); i++ {
		sum += args[i]
	}
	return sum
}

type QueryWrapper struct {
}

func (query *QueryWrapper) eq(s string, q string) *QueryWrapper {
	fmt.Printf("添加条件查询条件,%v == %v\n", s, q)
	return query
}

// 查询
func (q *QueryWrapper) query() {
	fmt.Println("执行sql查询")
}

func test(args ...any) { // 可以接受任意个string参数
	for _, v := range args {
		log.Println(v)
	}
}

func init() {
	log.Println("启动项目")
	// automatic.BuildReact(nil)
	go automatic.BuildJava(nil)
}

func main() {
	fmt.Println("测试 1")
	//定时任务更新梯子配置
	// trojan.GetYaml()
	// var identifier []trojan.MyServers
	// slice1 := make([]trojan.MyServers, 10)

	// for _, m := range trojan.GetConf().Proxies {
	// 	if strings.Contains(m["name"].(string), "日本") && m["type"] == "ss" {
	// 		// fmt.Printf("服务名: %v,密码:%v,端口:%v,类型:%v,服务地址:%v\n", m["name"], m["password"], m["port"], m["type"], m["server"])
	// 		slice1 = append(slice1, trojan.MyServers{
	// 			Address:  m["server"].(string),
	// 			Port:     m["port"].(int),
	// 			Method:   m["type"].(string),
	// 			Password: m["password"].(string),
	// 		})

	// 	}
	// }
	// fmt.Printf("slice1: %v\n", slice1)
	// trojan.GetJson(slice1)

	// r, _ := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
	// 	URL:      "https://github.com/go-git/go-git",
	// 	Progress: os.Stdout,
	// })
	// fmt.Printf("r: %v\n", r)
	// w, _ := r.Worktree()

	// w.Pull(&git.PullOptions{RemoteName: "origin"})
	// var arr = []string{
	//     "ele1",
	//     "ele2",
	//     "ele3",
	// }
	// test(arr[0],arr[1],arr[2])
	// test(flag.arr()...) // 切片被打散传入
	//链式调用封装测试
	// book := new(QueryWrapper)
	// book.eq("a","b").eq("c","k").query()

	// fmt.Printf("Sum([]int{1, 2, 3}...): %v\n", Sum([]id{1, 2, 3}...))
	// //接口测试
	// var _ Person = (*Worker)(nil)
	// stu := &Student{
	// 	name: "小明",
	// 	age: 14,
	// }
	// ex2(stu)
	// worker := &Worker{
	// 	name: "小李",
	// 	gender: "男",
	// }
	// //调用不同包路径下的函数
	// // fmt.Printf("calc.Add(1, 3): %v\n", calc.Add(1, 3))
	// //调用第三方包
	// // fmt.Printf("quote.Hello(): %v\n", quote.Hello())
	// ex(worker)

	// //
	// fmt.Printf("geeconfig.GetConf(): %v\n", geeconfig.GetConf())

	service := gee.NewService()

	service.POST("/pushReact", func(ctx *gee.Context) {
		go automatic.BuildReact(ctx)
	})

	service.POST("/pushBs", func(ctx *gee.Context) {
		go automatic.BuildJava(ctx)
	})
	service.Run(":9999")
}

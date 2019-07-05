package main

import (
	"fmt"
	"sort"
	"time"
)

//This ping function only accepts a channel for sending values. It would be a compile-time error to try to receive on this channel.
func ping(pings chan<- string, msg string) {
	pings <- msg
}
//The pong function accepts one channel for receives (pings) and a second for sends (pongs).
func pong(pings <-chan string, pongs chan<- string) {
	pongs <- <- pings
	//和以上等价 msg := <-pings
	//pongs <- msg
}

func main() {
	//多Goroutine值传递,并发打印输出
	names := []string{"Eric", "Harry", "Robert", "Jim", "Mark"}
	for _, name := range names {
		//fmt.Printf("Hello, %s.\n", name)
		//runtime.Gosched()
		go func(who string) {
			fmt.Printf("Hello, %s.\n", who)
		}(name)
	}
	time.Sleep(10 * time.Nanosecond)

	//Unbuffered channels 无缓冲通道
	//接收阻塞至收到消息，
	//发送阻塞直接收者接收到消息
	done := make(chan bool) //等价于 make(chan bool 0)
	nums := []int{2, 1, 3, 5, 4}
	go func() {
		time.Sleep(time.Second)
		sort.Ints(nums)
		done <- true
	}()
	<- done
	fmt.Println(nums)

	//Unbuffered channels 有缓冲通道
	//缓冲区未满时，发送者仅在值拷贝到缓冲区之前是阻塞的
	//而在缓冲区已满时，发送者会阻塞，直至接收者取走了消息，缓冲区有了空余
	messages := make(chan string, 2)//缓冲区大小为2
	messages <- "hello"
	messages <- "world"
	fmt.Println(<-messages, <-messages)

	//只读只写
	//ping函数发消息接收参数 约束为chan<-
	//pongs函数收消息接收参数分别为 接收与接收者 pings <-chan string, pongs chan<- string
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
package main

import (
	"fmt"
	"sort"
	"time"
)

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	pongs <- <-pings
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
	// 接收阻塞至收到消息，
	// 发送阻塞直接收者接收到消息
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
	// 缓冲区未满时，发送者仅在值拷贝到缓冲区之前是阻塞的
	//而在缓冲区已满时，发送者会阻塞，直至接收者取走了消息，缓冲区有了空余
	messages := make(chan string, 2)//缓冲区大小为2
	messages <- "hello"
	messages <- "world"
	fmt.Println(<-messages, <-messages)

    //函数封装时，对仅作消息接收或仅作消息发送的chan标识direction
    // 可以借用编译器检查增强类型使用安全。如下代码中，ping函数中pings chan仅用来接收消息，
    // 所以参数列表中将其标识为接收者。pong函数中，pings chan仅用来发送消息，
    // pongs chan仅用来接收消息，所以参数列表中二者分别标识为发送者与接收者。
	pings, pongs := make(chan string, 1), make(chan string, 1)
	ping(pings, "ping")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
package main

import (
	"errors"
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

func sum(nums []int) int {
	rlt := 0
	for _, num := range nums {
		rlt += num
	}
	return rlt
}

func sumWithTimeout(nums []int, timeout time.Duration) (int, error) {
	rlt := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		rlt <- sum(nums)
	}()
	select {
	case v := <-rlt:
		return v, nil
	case <-time.After(timeout): //以通道接收time.after定时器,select监听time.after定时器IO如果时间到达会返回
		return 0, errors.New("timeout")
	}
}

func main() {

    //1.Goroutine体验
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


	//2.1 无缓冲通道
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


    //2.2 有缓冲通道
	//Unbuffered channels 有缓冲通道
	//缓冲区未满时，发送者仅在值拷贝到缓冲区之前是阻塞的
	//而在缓冲区已满时，发送者会阻塞，直至接收者取走了消息,缓冲区有了空余.
	//注意缓冲区未满时非阻塞仅指发送者，对接收者依然是阻塞的
	messages := make(chan string, 2)//缓冲区大小为2
	messages <- "hello"
	messages <- "world"
	fmt.Println(<-messages, <-messages)


	//2.3通道 只读只写
	//ping函数发消息接收参数 约束为chan<-
	//pongs函数收消息接收参数分别为 接收与接收者 pings <-chan string, pongs chan<- string
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)


	//2.4 阻塞等待多个通道消息
	//select阻塞等待多个channel消息，select相当于监听io操作，事件产生后进行具体处理。
	//如下代码，创建两个chan，启动两个goroutine耗费不等时间计算结果，
	//主routine监听消息，使用两次select，
	//在for循环里，第一次进来接收到ch2消息，第二次接收到ch1消息，用时2.000521146s。
	c1, c2 := make(chan int, 0), make(chan int, 1)
	go func() {
		time.Sleep(2* time.Second)
		c1 <- 1
	}()
	go func() {
		time.Sleep(time.Second)
		c2 <- 2
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received msg from c1", msg1)
		case msg2 := <-c2:
			fmt.Println("received msg from c2", msg2)
		}
	}


	//2.5 select with default
	//select with default可以用来处理非阻塞式消息发送、接收及多路选择。
	//如下代码中第一个select为非阻塞式消息接收若收到消息，则落入<-messages case否则落入default。
	//第二个select为非阻塞式消息发送,与非阻塞式消息接收类似因messages chan为Unbuffered channel且无异步消息接收者，
	//因此落入default case。第三个select为多路非阻塞式消息接收
	messages2 := make(chan string)
	signal := make(chan bool)
	// receive with default
	for i := 0; i < 2; i++ {
		select {
		case <-messages2:
			fmt.Println("message received")
		default:
			fmt.Println("no message received")
		}
	}
	// send with default //这里为非阻塞式发送,message2为无缓存通道,下消息没有接收者读取是会阻塞的
	select {
	case messages2 <- "message":
		fmt.Println("message sent successfully")
	default:
		fmt.Println("message sent failed")
	}
	// muti-way select
	select {
	case <-messages2:
		fmt.Println("message received")
	case <-signal:
		fmt.Println("signal received")
	default:
		fmt.Println("no message or signal received")
	}


	//2.6 close
	//当无需再给channel发送消息时，可将其close类似进程通信的通道或socket文件描述符对端主动关闭。
	// 如下代码中创建一个Buffered channel，首先启动一个异步goroutine循环消费消息，
	//然后主routine完成消息发送后关闭chan，消费goroutine检测到chan关闭后，退出循环。
	messages3 := make(chan int, 10)
	done3 := make(chan bool)
	// consumer
	go func() {
		for {
			msg, more := <-messages3
			if !more {
				fmt.Println("no more message")
				done <- true
				break
			}
			fmt.Println("message received", msg)
		}
	}()
	// producer
	for i := 0; i < 5; i++ {
		messages3 <- i
	}
	close(messages3)
	<-done3


	//2.7 for range
	//for range语法不仅可对基础数据结构（slice、map等）作迭代，
	//还可对channel作消息接收迭代。如下代码中，给messages chan发送两条消息后将其关闭，
	//然后迭代messages chan打印消息。
	messages4 := make(chan string, 2)
	messages4 <- "hello"
	messages4 <- "world"
	close(messages4)

	for msg := range messages4 {
		fmt.Println(msg)
	}


	//2.8 超时控制
	//源访问、网络请求等场景作超时控制是非常必要的，
	//可以使用channel结合select来实现。如下代码，对常规sum函数增加超时限制，
	//sumWithTimeout函数中，select的v := <-rlt在等待计算结果，
	//若在时限范围内计算完成，则正常返回计算结果，
	//若超过时限则落入<-time.After(timeout) case，抛出timeout error。
	// 原理为，以通道接收time.after定时器,select监听time.after定时器IO如果时间到达会返回. 类似内核的alarm信号
	nums1 := []int{1, 2, 3, 4, 5}
	timeout := 3 * time.Second // time.Second
	rlt, err := sumWithTimeout(nums1, timeout)
	if nil != err {
		fmt.Println("error", err)
		return
	}
	fmt.Println(rlt)
}
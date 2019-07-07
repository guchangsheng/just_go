
//这是一个多协程同时操作共享变量出现竞争数据错误的问题
// This sample program demonstrates how to create race
// conditions in our programs. We don't want to do this.
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	// counter is a variable incremented by all goroutines.
	counter int
	// wg is used to wait for the program to finish.
	wg sync.WaitGroup

	mutex sync.Mutex
)

// main is the entry point for all Go programs.
func main() {
	//sysnc.waitGroup说明,
	//作用为主协程阻塞等待其它子协程完成工作。
	//流程为先通过add方法创建一个计数器，
	//再通过在协程执行方法中调用wg.Done()这里使用了refer 每调用一次计数器减1.
	//wg.Wait()等待直至计数器减少完成

	// Add a count of two, one for each goroutine.
	wg.Add(2)

	// Create two goroutines.
	go incCounter(1)
	go incCounter(2)

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

//程序分析 对共享全局变量进行非原子操作
//协程1 读取counter为0
//切换2
//协程2 读取counter为0
//切换1
//协程1 对counter++ counter变为1
//切换2
//协程2 对counter++ counter变为1
//切换1
//协程1 读取counter为1
//切换2
//协程2 读取counter为1
//切换1
//协程1 对counter++ counter变为2
//协程1 执行完成
//协程2 对counter++ counter变为3

//程序分析 增加互斥锁实现原子操作
//协程1 拿到互斥锁
//协程1 读取counter为0
//协程1 对counter++ counter变为1
//协程1 释放锁
//切换2
//协程2 拿到互斥锁
//协程2 读取counter为1
//协程2 对counter++ counter变为2
//协程2 释放锁
//切换1
//协程1 拿到互斥锁
//切换1 读取counter为2
//切换1 对counter++ counter变为3
//协程2 释放锁...

// incCounter increments the package level counter variable.
func incCounter(id int) {
	// Schedule the call to Done to tell main we are done.
	defer wg.Done()

	for count := 0; count < 2; count++{
		mutex.Lock()
		// Capture the value of Counter.
		value := counter
		// Yield the thread and be placed back in queue.
		runtime.Gosched() //主动让出cpu调度权限，让其它协程有被cpu调度机会
		// Increment our local value of Counter.
		value++
		// Store the value back into Counter.
		counter = value
		mutex.Unlock()
		fmt.Println("id:%d Counter:%d  time对counter++ counter变为2",id,counter,time.Now().UTC())
	}
}

package main

import (
	"fmt"
	"time"
)

func main() {
	firstDate := time.Now()
	defer func() {
		fmt.Println("耗时：", time.Since(firstDate))
	}()
	// 任务队列通道
	intChan := make(chan int, 1000)
	// 结果通道，所有计算的结果全部放在这里
	primeChan := make(chan int, 2000)
	// 标识退出的管道
	exitChan := make(chan int, 4)
	// 分发任务
	go putNum(intChan)
	// 开启四个协程来计算素数，并放入结果通道中
	for i := 0; i < 4; i++ {
		go cal(intChan, primeChan, exitChan)
	}
	// 开启协程，不断从exitChan中获取结束标志，当获取数量达到4个时
	go closeWork(primeChan, exitChan)
	// 主线程range遍历结果集
	for i := range primeChan {
		fmt.Println(i)
	}
	fmt.Println("遍历结束")
}

/*
*
PutNum：协程负责将所有需要计算的数字放入intChan通道
注意：全部放入后将intChan通道关闭，这样消费者通过for-range遍历时才不会死循环
*/
func putNum(intChan chan int) {
	for i := 1; i <= 10; i++ {
		intChan <- i
	}
	close(intChan)
}

/*
*
判断工作协程是否全部结束，如果结束则关闭primeChan，以此来通知主线程
*/
func closeWork(primeChan chan int, exitChan chan int) {
	for i := 0; i < 4; i++ {
		<-exitChan
	}
	close(primeChan)
	close(exitChan)
}

/*
*
for-range循环遍历intChan，并计算是否是素数。
for-range会遍历到该通道被关闭未知，当range循环结束后向exitChan中放入一个标识
表明当前协程已经结束
*/
func cal(intChan chan int, primeChan chan int, exitChan chan int) {
	for v := range intChan {
		flag := true
		for i := 2; i < v; i++ {
			if v%i == 0 {
				flag = false
				break
			}
		}
		if flag {
			primeChan <- v
		}
	}
	exitChan <- 0
}

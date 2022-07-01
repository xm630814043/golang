package channle

import (
	"fmt"
	"sync"
)

var (
	ch = make(chan int, 100)
)

func SendReads() {
	// 开始写入数据
	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch)
}

func ReceiveReads(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// 读取数据
	for i := range ch {
		fmt.Printf("receiver #%d get %d\n", id, i)
	}
	fmt.Printf("receiver #%d exit\n", id)
}

// SignalChannel 一写多读
//  唯一的写入端可以关闭 channel 用来通知读取端所有数据都已经写入完成了。读取端只需要用 for range 把 channel 中数据遍历完就可以了，
//  当 channel 关闭时，for range 仍然会将 channel 缓冲中的数据全部遍历完然后再退出循环
func SignalChannel() {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go ReceiveReads(0, wg)
	go ReceiveReads(1, wg)
	go ReceiveReads(2, wg)
	SendReads()
	wg.Wait()
}

// SendWrites 多写一读
// 解决多个写入端重复关闭 channel 的问题，但更优雅的办法设置一个额外的 channel ，由读取端通过关闭来通知写入端任务完成不要再继续再写入数据了
func SendWrites(id int, wg *sync.WaitGroup, done chan struct{}) {
	// 开始写入数据
	defer wg.Done()
	for i := 0; ; i++ {
		select {
		case <-done:
			// get exit signal
			fmt.Printf("sender #%d exit\n", id)
			return
		case ch <- id*1000 + i:
		}
	}
}

func ReceiveWrites(done chan struct{}) {
	count := 0
	for i := range ch {
		fmt.Printf("receiver get %d\n", i)
		count++
		if count >= 1000 {
			// signal recving finish
			close(done)
			return
		}
	}
}

func ChannelSignal() {
	wg := &sync.WaitGroup{}
	done := make(chan struct{})
	wg.Add(3)
	go SendWrites(0, wg, done)
	go SendWrites(1, wg, done)
	go SendWrites(2, wg, done)
	ReceiveWrites(done)
	wg.Wait()
}

package mutex

import (
	"fmt"
	"sync"
	"time"
)

// 多个线程在读相同的数据时
// 多个线程在写相同的数据时
// 同一个资源，有读又有写

var (
	count int
	lock  sync.Mutex
)

//互斥锁是一种常用的控制共享资源访问的方法，它能够保证同时只有一个 goroutine 可以访问到共享资源（同一个时刻只有一个线程能够拿到锁）
func Goroutine() {
	for i := 0; i < 2; i++ {
		go func() {
			for i := 1000000; i > 0; i-- {
				count++
			}
			fmt.Println(count)
		}()
	}
	fmt.Scanf("\n") //等待子线程全部结束
	fmt.Println("多线程:", count)
}

func read(mapA map[string]string) {
	for {
		lock.Lock()
		var _ string = mapA["name"]
		count += 1
		lock.Unlock()
	}
}

func write(mapA map[string]string) {
	for {
		lock.Lock()
		mapA["name"] = "johny"
		count += 1
		time.Sleep(time.Millisecond * 3)
		lock.Unlock()
	}
}

func Mutex() {
	for i := 0; i < 2; i++ {
		go func() {
			for i := 1000000; i > 0; i-- {
				lock.Lock()
				count++
				lock.Unlock()
			}
			fmt.Println(count)
		}()
	}
	fmt.Scanf("\n", count) //等待子线程全部结束
	fmt.Println("线程加互斥锁:", count)
}

func RunMutex() {
	var num int = 3
	var mapA map[string]string = map[string]string{"nema": ""}

	for i := 0; i < num; i++ {
		go read(mapA)
	}

	for i := 0; i < num; i++ {
		go write(mapA)
	}

	time.Sleep(time.Second * 3)
	fmt.Printf("互斥锁，最终读写次数：%d\n", count)
}

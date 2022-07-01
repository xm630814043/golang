package rwMutex

import (
	"fmt"
	"sync"
	"time"
)

// 在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效
// 如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
// 如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁

var (
	count  int
	rwLock sync.RWMutex
)

func read(mapA map[string]string) {
	for {
		rwLock.Lock()
		var _ string = mapA["name"]
		count += 1
		rwLock.Unlock()
	}
}

func write(mapA map[string]string) {
	for {
		rwLock.Lock()
		mapA["name"] = "johny"
		count += 1
		time.Sleep(time.Millisecond * 3)
		rwLock.Unlock()
	}
}

func RWMutex() {
	for i := 0; i < 2; i++ {
		go func() {
			for i := 1000000; i > 0; i-- {
				rwLock.Lock()
				count++
				rwLock.Unlock()
			}
			fmt.Println(count)
		}()
	}
	fmt.Scanf("\n", count) //等待子线程全部结束
	fmt.Println("读写锁")
}

func RunRWMutex() {
	var num int = 3
	var mapA map[string]string = map[string]string{"nema": ""}

	for i := 0; i < num; i++ {
		go read(mapA)
	}

	for i := 0; i < num; i++ {
		go write(mapA)
	}

	time.Sleep(time.Second * 3)
	fmt.Printf("读写锁，最终读写次数：%d\n", count)
}

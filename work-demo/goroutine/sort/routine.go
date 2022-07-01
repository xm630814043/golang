package sort

import (
	"fmt"
	"sync"
)

func TestGoroutineA() {
	array := []int{1, 4, 7, 6}

	var wg sync.WaitGroup
	for k, v := range array {
		wg.Add(1)
		fmt.Println("第", k, "次,A 的值: ", v)
		go func() {
			fmt.Println("协程打印输出", v)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestGoroutineB() {
	array := []int{1, 4, 7, 6}

	var wg sync.WaitGroup
	for k, v := range array {
		wg.Add(1)
		fmt.Println("第", k, "次,B 的值: ", v)
		go func(v int) {
			fmt.Println("协程打印输出", v)
			wg.Done()
		}(v)
	}
	wg.Wait()
}

func TestGoroutineC() {
	array := []int{1, 4, 7, 6}

	var wg sync.WaitGroup
	for k, v := range array {
		wg.Add(1)
		num := v
		fmt.Println("第", k, "次,C 的值: ", num)
		go func() {
			fmt.Println("协程打印输出", num)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test() {
	TestGoroutineA()
	TestGoroutineB()
	TestGoroutineC()
}

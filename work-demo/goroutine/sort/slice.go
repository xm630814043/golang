package sort

import (
	"fmt"
	"sync"
)

func SliceSort() {
	arr := []int{1, 4, 7}
	arr2 := []int{2, 5, 8}
	arr3 := []int{3, 6, 9}

	signalCh := make(chan int)
	signalCh2 := make(chan int)
	signalCh3 := make(chan int)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range arr {
			fmt.Println(v)
			signalCh <- 1
			<-signalCh3
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range arr2 {
			<-signalCh
			fmt.Println(v)
			signalCh2 <- 1
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range arr3 {
			<-signalCh2
			fmt.Println(v)
			signalCh3 <- 1
		}
	}()
	wg.Wait()
}

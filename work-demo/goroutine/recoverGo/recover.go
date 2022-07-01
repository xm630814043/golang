package recoverGo

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func NewGo(f func()) {
	_, file, line, ok := runtime.Caller(1)
	go func() {
		defer func() {
			err := recover()
			if ok {
				log.Printf("协程发生未知错误:path:%s,line:%d,Error:%+v", file, line, err)
			} else {
				log.Printf("协程成功", file, line)
			}
		}()
		f()
	}()
}

func DemoWrite() {
	f, err := os.Create("output.txt") // creating...
	if err != nil {
		fmt.Printf("error creating file: %v", err)
		return
	}
	defer f.Close()
	for i := 1; i < 100; i++ { // Generating...
		_, err = f.WriteString(fmt.Sprintf("%d|", i)) // writing...
		if err != nil {
			fmt.Printf("error writing string: %v", err)
		}
	}
}

func DemoGo() {
	arr := []int{1}
	a := func() {
		iserr := arr[2]
		log.Printf("go协程输出：", iserr)
	}
	NewGo(a)
	time.Sleep(time.Second * 2)
	log.Printf("最终输出：", arr)
}

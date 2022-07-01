package copy

import (
	"fmt"
	"log"
)

func Demo() {
	//数组
	arrayb := [6]int{10, 20, 30, 40}
	sliceb := [6]int{6}
	fmt.Println("数组：", arrayb)
	fmt.Println("数组：", sliceb)
	//切片
	arraya := []int{10, 20, 30, 40}
	slicea := make([]int, 6)
	num := copy(slicea, arraya)
	fmt.Println("切片：", num)
	fmt.Println("切片：", slicea)
	log.Fatalln(fmt.Sprintf("num", slicea))
	fmt.Sprintf("num", slicea)
}

package main

import (
	"fmt"
	"work-demo/agent/snmp"
)

func SplitArray(array []int, part int) {
	arrayList := make([][]int, part)
	num := len(array) / part
	var j, k int
	for i := 0; i < len(array); i++ {
		if k == part-1 {
			arrayList[k] = append(arrayList[k], array[i])
			continue
		}
		if j == num {
			j = 0
			k++
		}
		arrayList[k] = append(arrayList[k], array[i])
		j++
	}
	fmt.Println(arrayList)
}

func deferCall() {
	defer func() { fmt.Println("打印前") }()
	defer func() { fmt.Println("打印中") }()
	defer func() { fmt.Println("打印后") }()

	panic("触发异常")
}

func main() {
	//sort.TestChan()
	//sort.SliceSort()
	//sort.ForSort()
	//获取当前系统
	//fmt.Println(agent.Ping("www.baidu.com", time.Second))

	//_ = feishu.FSErrorNotify(WebhookUrl, "数据采集", errors.New("11111111"), "")
	//
	//sort.Test()
	//agent.Ping("www.baidu.com", time.Second)

	//consumer.RunConsumer()
	//producer.RunProducer()
	//
	//demo.Demo()
	snmp.TrapAccept()

}

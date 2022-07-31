package main

import (
	"fmt"
	"work-demo/agent"
)

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

	// MQ
	//consumer.RunConsumer()
	//producer.RunProducer()
	//consumer.RunConsumers()
	//producer.RunProducers()

	//demo.Demo()
	//snmp.TrapAccept()
	ip, _ := agent.ExternalIP()
	fmt.Printf(ip.String())
}

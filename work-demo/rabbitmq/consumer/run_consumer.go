package consumer

import (
	"fmt"
)

type TestPro struct {
	msgContent string
}

// 实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}

func RunConsumer() {
	msg := fmt.Sprintf("这是消费者测试任务")
	t := &TestPro{
		msg,
	}
	queueExchange := &QueueExchange{
		"test_rabbit",
		"",
		"test_rabbit",
		"direct",
	}
	mq := New(queueExchange.QuName, queueExchange.ExName, queueExchange.ExType, "")
	//mq.RegisterReceiver(t)
	mq.listenReceiver(t)
	mq.Start()
	fmt.Println("Consumer：消费者消费成功")
}

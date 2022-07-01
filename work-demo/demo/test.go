package demo

import (
	"encoding/json"
	"fmt"
)

// 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	url          string //MQ链接字符串
	queueName    string // 队列名称
	routingKey   string // key名称
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
}

type Test struct {
	QueueExchange []*QueueExchange `json:"queue_exchange"`
	RabbitMQ      []*RabbitMQ      `json:"rabbit_mq"`
}

func Demo() {
	rabbitMQs := make([]*RabbitMQ, 0)
	rabbitMQs = append(rabbitMQs, &RabbitMQ{
		url:       "111",
		queueName: "222",
	})
	exchanges := make([]*QueueExchange, 0)
	exchanges = append(exchanges, &QueueExchange{
		QuName: "333",
		RtKey:  "444",
	})

	test := Test{
		QueueExchange: exchanges,
		RabbitMQ:      rabbitMQs,
	}

	marshal, _ := json.Marshal(test)
	//fmt.Printf("marshal: %v \n", marshal)
	tests := &Test{}
	_ = json.Unmarshal(marshal, &tests)
	fmt.Printf("Unmarshal: %v \n", tests.QueueExchange[0].QuName)
}

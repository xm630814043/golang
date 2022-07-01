package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const Mqurl = "amqp://admin:123456@49.234.72.129:5671/"

var (
	// 定义全局变量,指针类型
	mqConn *amqp.Connection
	mqChan *amqp.Channel
)

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	url          string //MQ链接字符串
	conn         *amqp.Connection
	channel      *amqp.Channel
	queueName    string // 队列名称
	routingKey   string // key名称
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	receiverList []Receiver
	mu           sync.RWMutex
}

// 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

// 创建一个新的操作对象
func New(queueName, exchangeName, exchangeType, routingKey string) *RabbitMQ {
	rabbitMQ := RabbitMQ{
		queueName:    queueName,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		routingKey:   routingKey,
		url:          Mqurl,
	}
	return &rabbitMQ
}

// 链接rabbitMQ
func (r *RabbitMQ) mqConnect() {
	var err error
	time.Sleep(2 * time.Second)

	mqConn, err = amqp.Dial(r.url)
	fmt.Println("rabbitMQ链接地址：", r.url)
	r.conn = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开链接失败:%s \n", err)
	}
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
	}
}

// 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.mu.Lock()
	r.receiverList = append(r.receiverList, receiver)
	r.mu.Unlock()
}

// 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) Start() {
	// 开启监听消费者发送任务
	for _, receiver := range r.receiverList {
		fmt.Println("开启监听消费者发送任务", receiver)
		go r.listenReceiver(receiver)
	}
	time.Sleep(1 * time.Second)
}

// 关闭RabbitMQ连接,释放资源,建议NewRabbitMQ获取实例后
func (r *RabbitMQ) mqClose() {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = r.conn.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver) {
	// 处理结束关闭链接
	defer r.mqClose()
	// 验证链接是否正常
	if r.channel == nil {
		r.mqConnect()
	}
	// 长连接：eventbasicconsumer + noack.... 【订阅式】,consumer端处理一条数据需要耗费 1s钟。。。。
	//《1》 确认机制。。。 不管你是否却不确认，消息都会一股脑全部打入到你的consumer中去。。。
	//《2》 QOS =》 服务质量。。。 【QOS + Ack】机制，解决这个问题。。。
	//解决办法就是在channel设置好通道。。。
	//channel.Qos 这样RabbitMQ就会使得每个Consumer在同一个时间点最多处理一个Message。换句话说，在接收到该Consumer的ack前，他它不会将新的Message分发给它。
	//param1：prefetchSize，预取大小服务器将传递的最大内容量（以八位字节为单位），如果不受限制，则为0;默认值：0
	//param2：prefetchCount，服务器一次请求将传递的最大邮件数，如果没有限制，则为0;调用此方法时，该值必填。默认值：0
	//param3：global，是否将设置应用于整个频道，而不是每个消费者;默认值：false，应用于本身（一个消费者）,true：应用于整个频道
	_ = r.channel.Qos(30, 0, false)

	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.channel.QueueDeclarePassive(r.queueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 绑定任务
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)
	msgList, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
		return
	}
	for msg := range msgList {
		// 处理数据
		fmt.Println("处理数据:", msg.Body)
		err := receiver.Consumer(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			fmt.Println("消费失败,消息体数据：", string(msg.Body))
			continue // 消费失败 消息重新放回队列
		} else {
			// 确认消息,必须为false
			_ = msg.Ack(true)
			fmt.Println("消费成功；消息体数据：", string(msg.Body))
			return
		}
	}
}

package producer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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

// 定义生产者接口
type Producer interface {
	MsgContent() string
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
	producerList []Producer
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

	cfg := new(tls.Config)
	cfg.RootCAs = x509.NewCertPool()
	// 加载机构颁发证书
	if ca, err := ioutil.ReadFile("rabbitmq/testca/cacert.pem"); err == nil {
		cfg.RootCAs.AppendCertsFromPEM(ca)
	}
	// 加载客户端证书和密钥
	if cert, err := tls.LoadX509KeyPair("rabbitmq/client/cert.pem", "rabbitmq/client/key.pem"); err == nil {
		cfg.Certificates = append(cfg.Certificates, cert)
	}
	fmt.Println("MQ请求链接：", r.url)

	mqConn, err = amqp.DialTLS(r.url, cfg)
	if err != nil {
		fmt.Printf("MQ打开链接失败conn: %v, err: %v \n", mqConn, err)
	}
	r.conn = mqConn // 赋值给RabbitMQ对象
	mqChan, err = mqConn.Channel()
	r.channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
	}

}

// 启动RabbitMQ客户端,并初始化
func (r *RabbitMQ) Start() {
	// 开启监听生产者发送任务
	for _, producer := range r.producerList {
		//fmt.Println("开启监听生产者发送任务", producer)
		go r.listenProducer(producer)
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

// 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.producerList = append(r.producerList, producer)
}

// 发送任务
func (r *RabbitMQ) listenProducer(producer Producer) {
	defer r.mqClose()
	// 验证链接是否正常,否则重新链接
	if r.channel == nil {
		r.mqConnect()
	}
	fmt.Println("验证链接正常,发送任务", r)
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
	// 队列绑定
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return
	}
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = r.channel.ExchangeDeclarePassive(r.exchangeName, r.routingKey, true, false, false, true, nil)
	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = r.channel.ExchangeDeclare(r.exchangeName, r.routingKey, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return
		}
	}
	// 发送任务消息
	err = r.channel.Publish(r.exchangeName, r.routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
	})
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

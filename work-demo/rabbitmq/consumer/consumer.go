package consumer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const Mqurl = "amqps://admin:123456@49.234.72.129:5671/"

var (
	// 定义全局变量,指针类型
	rabbitmqConn *amqp.Connection
	rabbitmqChan *amqp.Channel
)

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

// 定义RabbitMQ对象
type RabbitMQ struct {
	url          string //MQ链接字符串
	connection   *amqp.Connection
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
func (r *RabbitMQ) mqConnect() error {
	// 加载机构颁发证书
	caCert, err := ioutil.ReadFile("rabbitmq/cacert/cacert.pem")
	if err != nil {
		log.Fatal("加载机构颁发证书: ", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 加载客户端证书和密钥
	cert, err := tls.LoadX509KeyPair("rabbitmq/client/rabbit-client.cert.pem", "rabbitmq/client/rabbit-client.key.pem")
	if err != nil {
		log.Fatal("加载客户端证书和密钥: ", err)
		return err
	}
	fmt.Println("MQ请求链接：", r.url)

	rabbitmqConn, err = amqp.DialTLS(r.url, &tls.Config{
		Certificates:             []tls.Certificate{cert}, // from tls.LoadX509KeyPair
		RootCAs:                  caCertPool,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		InsecureSkipVerify:       true,
		ServerName:               "VM-0-10-centos",
		MinVersion:               tls.VersionTLS10,
	})
	if err != nil {
		fmt.Println("MQ打开链接失败: ", err)
		return err
	}

	fmt.Println("MQ链接请求返回值mqConn：", rabbitmqConn)

	r.connection = rabbitmqConn // 赋值给RabbitMQ对象
	rabbitmqChan, err = rabbitmqConn.Channel()
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
		return err
	}
	fmt.Println("MQ打开管道返回值mqChan：", rabbitmqChan)
	r.channel = rabbitmqChan // 赋值给RabbitMQ对象
	return nil
}

// 关闭RabbitMQ连接
func (r *RabbitMQ) mqClose() error {
	// 先关闭管道,再关闭链接
	err := r.channel.Close()
	if err != nil {
		log.Fatal("关闭管道err: ", err)
		return err
	}
	// 关闭mq连接
	err = r.connection.Close()
	if err != nil {
		log.Fatal("关闭链接err: ", err)
		return err
	}
	return nil
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
		go r.listenReceiver(receiver)
	}
	time.Sleep(1 * time.Second)
}

// 监听接收者接收任务
func (r *RabbitMQ) listenReceiver(receiver Receiver) {
	var err error
	err = r.mqConnect()
	if err != nil {
		fmt.Println("启动RabbitMQ客户端,并初始化失败: ", err)
	}
	// 验证链接是否正常,否则重新链接
	if r.channel == nil || r.connection.IsClosed() {
		err = r.mqConnect()
		if err != nil {
			fmt.Println("验证链接是否正常,否则重新链接,初始化失败: ", err)
			return
		}
	}
	fmt.Println("验证链接正常,发送任务", r)
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
	_, err = r.channel.QueueDeclarePassive(r.queueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.channel.QueueDeclare(r.queueName, true, false, false, true, nil)
		if err != nil {
			fmt.Println("MQ注册队列失败: ", err)
			return
		}
	}
	// 绑定任务
	err = r.channel.QueueBind(r.queueName, r.routingKey, r.exchangeName, true, nil)
	if err != nil {
		fmt.Println("绑定队列失败: ", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.channel.Qos(1, 0, true)
	msgList, err := r.channel.Consume(r.queueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("获取消费通道异常: ", err)
		return
	}
	for msg := range msgList {
		// 处理数据
		//fmt.Println("处理数据:", msg.Body)
		err := receiver.Consumer(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			fmt.Println("消费失败,消息体数据：", string(msg.Body))
			continue // 消费失败 消息重新放回队列
		} else {
			// 确认消息,必须为false
			_ = msg.Ack(true)
			fmt.Println("消费成功；消息体数据：", string(msg.Body))
		}
	}
	err = r.mqClose()
	return
}

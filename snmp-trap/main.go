package main

import (
    "encoding/json"
    "fmt"
    "github.com/imroc/req"
    "github.com/pkg/errors"
    "net"
    "net/http"
    "time"
    "work-demo/feishu"
    "work-demo/snmp"
)

type MetricValue struct {
    Endpoint string      `json:"endpoint"`    //标明Metric的主体(属主)，比如metric是cpu_idle，那么Endpoint就表示这是哪台机器的cpu_idle
    Metric   string      `json:"metric"`      //最核心的字段，代表这个采集项具体度量的是什么, 比如是cpu_idle呢，还是memory_free, 还是qps
    Value    interface{} `json:"value"`       //代表该metric在当前时间点的值，float64
    Step     int64       `json:"step"`        //表示该数据采集项的汇报周期，这对于后续的配置监控策略很重要，必须明确指定
    Type     string      `json:"counterType"` //只能是COUNTER或者GAUGE二选一，前者表示该数据采集项为计时器类型，后者表示其为原值 (注意大小写)；
    // GAUGE：即用户上传什么样的值，就原封不动的存储。COUNTER：指标在存储和展现的时候，会被计算为speed，即（当前值 - 上次值）/ 时间间隔
    Tags      string `json:"tags"`      //一组逗号分割的键值对, 对metric进一步描述和细化, 可以是空字符串. 比如idc=lg，比如service=xbox等，多个tag之间用逗号分割
    Timestamp int64  `json:"timestamp"` //表示汇报该数据时的unix时间戳，注意是整数，代表的是秒
}

type TrapMetricValue struct {
    MetricValue []*MetricValue   `json:"metric_value"`
    SnmpPacket  *snmp.SnmpPacket `json:"packet"`
    Addr        *net.UDPAddr     `json:"addr"`
}

func myTrapHandler(packet *snmp.SnmpPacket, addr *net.UDPAddr) {

    // MetricValue 结构体组装
    arrayMetricValue := make([]*MetricValue, 0)
    for _, v := range packet.Variables {
        fmt.Printf("接受到的TRAP信息:%s \n", v)
        switch v.Type {
        case snmp.OctetString:
            b := v.Value.([]byte)
            arrayMetricValue = append(arrayMetricValue, &MetricValue{
                Endpoint:  v.Name + "_" + "test_shuai" + "_" + "5016fc23e6f6424fbb22135f42402bbb",
                Metric:    "OID",
                Value:     b,
                Step:      60,
                Type:      "GAUGE",
                Tags:      "",
                Timestamp: time.Now().Unix(),
            })
        default:
            arrayMetricValue = append(arrayMetricValue, &MetricValue{
                Endpoint:  v.Name + "_" + "test_shuai" + "_" + "5016fc23e6f6424fbb22135f42402bbb",
                Metric:    v.TypeName,
                Value:     v.Value,
                Step:      60,
                Type:      "GAUGE",
                Tags:      "",
                Timestamp: time.Now().Unix(),
            })
        }
    }
    fmt.Printf("接受到的TRAP信息:%s \n", packet)
    // 转发Transfer
    TrapSendTransfer(arrayMetricValue)

    // 飞书日志打印
    _ = feishu.FSErrorNotify(feishu.WebhookUrl, "trap监听消息接受成功", nil, TrapMetricValue{
        MetricValue: arrayMetricValue,
        SnmpPacket:  packet,
        Addr:        addr,
    })
}

func TrapSendTransfer(metric []*MetricValue) {
    // url := global.GlobalCfg.DataServiceIp + ":" + strconv.Itoa(global.GlobalCfg.DataServicePort) + "/api/push"
    url := "http://172.31.0.151:6060" + "/api/push"
    header := make(http.Header)
    header.Set("Content-Type", "application/json; charset=UTF-8")
    data, _ := json.Marshal(metric)

    for {
        r, err := req.Post(url, header, data)
        if r != nil {
            if err != nil || r.Response().StatusCode != 200 {
                return
            }
            break
        }
        time.Sleep(3 * time.Second)
    }
}

func main() {
    defer func() {
        if rec := recover(); rec != nil {
            _ = feishu.FSErrorNotify(feishu.WebhookUrl, "trap监听消息接受失败", errors.New(fmt.Sprintf("%v", rec)), "")
        }
    }()
    fmt.Println("进入trap接受")
    tl := snmp.NewTrapListener()
    tl.OnNewTrap = myTrapHandler
    tl.Params = snmp.Default

    err := tl.Listen("0.0.0.0:162")
    if err != nil {
        _ = feishu.FSErrorNotify(feishu.WebhookUrl, "trap监听消息接受错误", err, tl)
    }

}

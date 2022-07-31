package trap

import (
    "encoding/json"
    "fmt"
    "github.com/deejross/go-snmplib"
    "github.com/imroc/req"
    log "github.com/sirupsen/logrus"
    "go-snmplib/feishu"
    "math/rand"
    "net"
    "net/http"
    "time"
)

var udpSockTrap *net.UDPConn
var snmpTrap *snmplib.SNMP

func InitMQ() {
    log.Info("TRAP监听业务代码初始化")
    StartListenTrapMessage()
}

type MetricValue struct {
    Endpoint  string      `json:"endpoint"`
    Metric    string      `json:"metric"`
    Value     interface{} `json:"value"`
    Step      int64       `json:"step"`
    Type      string      `json:"counterType"`
    Tags      string      `json:"tags"`
    Timestamp int64       `json:"timestamp"`
}

func StartUDPServer(listenIPAddr string, port int) *net.UDPConn {
    log.Info("启动 udp listen")
    addr := net.UDPAddr{
        Port: port,
        IP:   net.ParseIP(listenIPAddr),
    }
    conn, err := net.ListenUDP("udp", &addr)
    if err != nil {
        log.Warn("udp listen error ", err.Error())
        _ = feishu.FSErrorNotify(feishu.FeiShuWebhookUrl, "trap启动UDP失败", err, "")
    }
    return conn
}

// snmptrap -v 3 -a SHA -A trap_net_works -x AES -X net_works_trap -l authPriv -u pcb.snmpv3 -e 0x80001f8880315de44d53ce8394 127.0.0.1 "" linkUp.0 1.3.6.4.5.6 s "switch is off"
// snmptrap -v 2c -c public 127.0.0.1 '' SNMPv2-MIB::system SNMPv2-MIB::sysDescr.0 s "red laptop" SNMPv2-MIB::sysServices.0 i "5" SNMPv2-MIB::sysObjectID o "1.3.6.1.4.1.2.3.4.5"

// @Title StartListenTrapMessage
// @Description  trap启动监听
// @Author zhaopengtao 2022-07-06 16:01:47
// @Param ctxMain
// @Param cfg

func StartListenTrapMessage() {
    rand.Seed(0)
    // 启动监听
    udpSockTrap = StartUDPServer("0.0.0.0", 162)

    // Target 目标设备，Community 设备的地址
    log.Info("启动 snmp client")
    snmpTrap = snmplib.NewSNMPOnConn("", "", snmplib.SNMPv3, 2*time.Second, 5, udpSockTrap)
    //snmp := snmplib.NewSNMPOnConn("", "", snmplib.SNMPv2c, 2*time.Second, 5, udpSock)

    // v3协议
    snmpTrap.TrapUsers = append(snmpTrap.TrapUsers, snmplib.V3user{"pcb.snmpv3", "SHA1", "trap_net_works", "AES", "net_works_trap"})

    // 获取监听到的报文数据

    packet := make([]byte, 3000)
    _, _, err := udpSockTrap.ReadFromUDP(packet)
    if err != nil {
        log.Warn("udp read error", err.Error())
        _ = feishu.FSErrorNotify(feishu.FeiShuWebhookUrl, "tudp read error", err, udpSockTrap)
    }
    //log.Info("Received trap from : ", addr.IP)
    //addrIP := fmt.Sprintf("%s", addr.IP)
    // 解析监听到的报文数据
    varBinds, err := snmpTrap.ParseTrap(packet)
    if err != nil {
        log.Warn("Error processing trap:", err.Error())
        _ = feishu.FSErrorNotify(feishu.FeiShuWebhookUrl, "Error processing trap", err, packet)
    }
    fmt.Sprintf("%v \n", varBinds)
    // 转发Transfer
    //todo 全局配置没有生成，暂不打开
    //SetMetricValueSendTransfer(addrIP, varBinds.VarBindOIDs)

}

func StopListenTrapMessage() {
    log.Info("close listen trap")
    udpSockTrap.Close()
    snmpTrap.Close()
}

// @Title SetMetricValueSendTransfer
// @Description  整理转发Transfer结构体并发送
// @Author zhaopengtao 2022-07-06 14:19:22
// @Param listenIPAddr
// @Param port
// @Return *net.UDPConn

func SetMetricValueSendTransfer(addrIP string, oIds []string) {
    log.Info("进入整合Transfer结构体,转发Transfer")
    metricValues := make([]*MetricValue, 0)
    // 整合Transfer结构体
    metricValue := &MetricValue{
        Endpoint:  "listenTrap:" + addrIP + "_" + "test_shuai" + "_" + "5016fc23e6f6424fbb22135f42402lll",
        Metric:    "OID",
        Value:     oIds,
        Step:      60,
        Type:      "GAUGE",
        Tags:      "",
        Timestamp: time.Now().Unix(),
    }
    // 转发Transfer
    metricValues = append(metricValues, metricValue)
    TrapSendTransfer(metricValues)
    log.Info("成功转发Transfer........  ")
    // 飞书日志打印
    _ = feishu.FSErrorNotify(feishu.FeiShuWebhookUrl, "trap监听消息成功", nil, metricValue)
}

func TrapSendTransfer(metric []*MetricValue) {
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
        log.Info("转发Transfer中........  ")
        time.Sleep(3 * time.Second)
    }
}

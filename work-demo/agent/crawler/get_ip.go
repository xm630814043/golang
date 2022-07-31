package crawler

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func GetIpAddress() {
	var err error
	var arrIpAddr []string // 本地IP

	arrIpAddr, err = getLocalIPv4s()
	if nil != err {
		fmt.Printf("err= %v \n", err)
	}
	fmt.Println("本地ip ", strings.Join(arrIpAddr, ","))

	var systemType, system, hostname string
	systemType = runtime.GOOS
	system = runtime.GOARCH
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Printf("err= %v \n", err)
	}
	fmt.Println("当前系统：", systemType+"-"+system+"-"+hostname)
}

// getLocalIPv4s 通过 net.InterfaceAddrs 获取本地IP
func getLocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
		}
	}

	return ips, nil
}

// getExternalIp 百度输入框中输入IP，即可查看公网ip
func getExternalIp() string {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		fmt.Println("get external ip err=", err)
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

// IsPublicIP 判断是否是公网IP
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

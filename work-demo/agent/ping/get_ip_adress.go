package ping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data"`
}

type IP struct {
	Country   string `json:"country"`
	CountryId string `json:"country_id"`
	Area      string `json:"area"`
	AreaId    string `json:"area_id"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	City      string `json:"city"`
	CityId    string `json:"city_id"`
	Isp       string `json:"isp"`
}

// TabaoAPI 通过淘宝接口根据公网ip获取国家运营商等信息
func TabaoAPI(ip string) *IPInfo {
	url := "http://ip.taobao.com/service/getIpInfo.php?ip="
	url += ip
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return &result
}

func GetIpAddress() {
	//arrIpAddr, err := getLocalIPv4s()
	//if nil != err {
	//    fmt.Printf("err=%v", err)
	//} else {
	//    for _, ip := range arrIpAddr {
	//        fmt.Println("get 本地 ip", ip)
	//        result := TabaoAPI(ip)
	//        if result != nil {
	//            fmt.Println("国家：", result.Data.Country)
	//            fmt.Println("地区：", result.Data.Area)
	//            fmt.Println("城市：", result.Data.City)
	//            fmt.Println("运营商：", result.Data.Isp)
	//        }
	//    }
	//}

	arrIpAddr, err := getLocalIPv4s()
	if nil != err {
		fmt.Printf("err=%v", err)
	}
	fmt.Println("get 本地 ip", arrIpAddr)
	strExternalIp := getExternalIp()
	isPublicIpAddr := IsPublicIP(net.ParseIP(strExternalIp))
	if isPublicIpAddr {
		fmt.Println("是公网ip ", strExternalIp)
	} else {
		fmt.Println("否公网ip ", strExternalIp)
	}

	//result := TabaoAPI(strExternalIp)
	//if result != nil {
	//    fmt.Println("国家：", result.Data.Country)
	//    fmt.Println("地区：", result.Data.Area)
	//    fmt.Println("城市：", result.Data.City)
	//    fmt.Println("运营商：", result.Data.Isp)
	//}
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

func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func inet_ntoa(ipnr int64) net.IP {
	var byteSlice [4]byte
	byteSlice[0] = byte(ipnr & 0xFF)
	byteSlice[1] = byte((ipnr >> 8) & 0xFF)
	byteSlice[2] = byte((ipnr >> 16) & 0xFF)
	byteSlice[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(byteSlice[3], byteSlice[2], byteSlice[1], byteSlice[0])
}

// IpBetween 判断ip地址区间
func IpBetween(from net.IP, to net.IP, test net.IP) bool {
	if from == nil || to == nil || test == nil {
		fmt.Println("An ip input is nil")
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := test.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		return false
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true
	}
	return false
}

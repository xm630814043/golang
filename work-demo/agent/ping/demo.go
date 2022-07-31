package ping

import (
	"fmt"
	"github.com/ip2location/ip2location-go"
)

func ReadableSize() {
	//(文件：IP2LOCATION-LITE-DB3.IPV6.BIN) 下载地址：https://lite.ip2location.com/
	db, err := ip2location.OpenDB("./IP2LOCATION-LITE-DB5.BIN")

	if err != nil {
		fmt.Print(err)
		return
	}
	ip := "58.37.2.34"
	results, err := db.Get_all(ip)

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Printf("ip: %s\n", ip)
	fmt.Printf("country: %s\n", results.Country_long)
	fmt.Printf("region: %s\n", results.Region)
	fmt.Printf("city: %s\n", results.City)
	fmt.Printf("latitude: %f\n", results.Latitude)
	fmt.Printf("longitude: %f\n", results.Longitude)

	db.Close()
}

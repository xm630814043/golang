package agent

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func Ping(addr string, t time.Duration) bool {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return false
	}
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	st := pinger.Statistics()
	fmt.Printf("[%.v]", st.PacketLoss)
	fmt.Printf("\n--- %s ping statistics ---\n", st.Addr)

	fmt.Printf("ip %s, %d packets transmitted, %d packets received, %v%% packet loss\n",

		st.PacketsSent, st.PacketsRecv, st.PacketLoss)

	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",

		st.MinRtt, st.AvgRtt, st.MaxRtt, st.StdDevRtt)

	/*
		设置pinger将发送的类型。
		false表示pinger将发送“未经授权”的UDP ping
		true表示pinger将发送“特权”原始ICMP ping
	*/
	//fmt.Println(pinger)
	//fmt.Sprintf("%v", pinger)

	pinger.SetPrivileged(true)
	// 运行pinger
	pinger.Run()
	stats := pinger.Statistics()

	//fmt.Sprintf("%v", stats)
	if stats.PacketsRecv >= 1 {
		return true
	} else {
		return false
	}
}

// +build linux

package bandwidth

import (
	"fmt"
	"log"
	"time"

	"github.com/juanhuttemann/nitr-api/network"
	"github.com/prometheus/procfs"
)

type NetworkDeviceBandwidth struct {
	Name      string `json:"name"`
	RxBytes   uint64 `json:"rxBytes"`
	TxBytes   uint64 `json:"txBytes"`
	RxPackets uint64 `json:"rxPackets"`
	TxPackets uint64 `json:"txPackets"`
}

func Check() []IfaceStats {
	p, err := procfs.NewDefaultFS()
	if err != nil {
		log.Fatalf("could not get process: %s", err)
	}
	net, err := p.NetDev()
	if err != nil {
		fmt.Println(err)
	}

	networks := network.Check()

	//Round 1
	var stats1 []IfaceStats
	for _, netw := range networks {
		stats1 = append(stats1, IfaceStats{
			Name:      netw.Name,
			RxBytes:   net[netw.Name].RxBytes,
			TxBytes:   net[netw.Name].TxBytes,
			RxPackets: net[netw.Name].RxPackets,
			TxPackets: net[netw.Name].TxPackets,
		})
	}

	time.Sleep(1000 * time.Millisecond)

	net, err = p.NetDev()
	if err != nil {
		fmt.Println(err)
	}

	//Round 2
	var stats2 []IfaceStats
	for _, netw := range networks {
		stats2 = append(stats2, IfaceStats{
			Name:      netw.Name,
			RxBytes:   net[netw.Name].RxBytes,
			TxBytes:   net[netw.Name].TxBytes,
			RxPackets: net[netw.Name].RxPackets,
			TxPackets: net[netw.Name].TxPackets,
		})
	}

	//DIFF
	var diffStats []IfaceStats

	for i, netw := range networks {
		diffStats = append(diffStats, IfaceStats{
			Name:      netw.Name,
			RxBytes:   stats2[i].RxBytes - stats1[i].RxBytes,
			TxBytes:   stats2[i].TxBytes - stats1[i].TxBytes,
			RxPackets: stats2[i].RxPackets,
			TxPackets: stats2[i].TxPackets,
		})
	}

	return diffStats
}

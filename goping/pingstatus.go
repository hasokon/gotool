package main

import (
	"fmt"
	)

type PingStatus struct {
	DataSize uint
	Hostname string
	IPAddr string
	Seq uint
	TTL uint
	RTT float64
}

func (ps PingStatus) String() string {
	return fmt.Sprintf("%d bytes from %s (%s) : icmp_seq=%d ttl=%d time=%.3f ms", ps.DataSize, ps.Hostname, ps.IPAddr, ps.Seq, ps.TTL, ps.RTT * 1000)
}

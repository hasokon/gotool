package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"math"
	"work/gotool/goping/icmp"
	"work/gotool/goping/ipheader"
)

func getIPAddr(host string) (net.IP, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.To4(), nil
		}
	}

	return nil, errors.New("IP address not found")
}

func getHostname(ip string) (string, error) {
	hostnames, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}

	if len(hostnames) > 0 {
		return hostnames[0], nil
	}

	return "", errors.New("Host Name not found")
}

func pinger(connection net.Conn, i *icmp.ICMP, exitNotifier chan int, sigint chan os.Signal) {
	seq := uint16(0)
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

	done := false
	for !done {
		select {
		case <-sigint:
			done = true
			break
		case <-t.C:
			timeBinary, err := time.Now().MarshalBinary()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			seq++ // This value also indicates the number of packets sent
			i.Seq = seq
			i.Data = timeBinary

			_, err = connection.Write(i.Marshal())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	}

	exitNotifier <- int(seq)
}

func getPingStatus(receiveData []byte, n int) PingStatus {
	// Parse IP Header
	iphdrLen := uint8(receiveData[0]&0xf) * 4
	iphdr := &ipheader.IPHeader{}
	err := iphdr.Parse(receiveData[:iphdrLen])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//Get IP Addr and Host Name
	dataSize := uint(iphdr.TotalLen)
	ttl := uint(iphdr.TTL)
	ipaddr := iphdr.SrcAddrString()
	hostname, err := getHostname(ipaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Parse ICMP Message
	receiveData = receiveData[iphdrLen:n]
	receiveMessage := &icmp.ICMP{}
	receiveMessage.ParseEchoMessage(receiveData)
	seq := uint(receiveMessage.Seq)

	transTime := time.Time{}
	transTime.UnmarshalBinary(receiveMessage.Data)
	rtt := time.Now().Sub(transTime)

	return PingStatus {
		DataSize : dataSize,
		Hostname : hostname,
		IPAddr : ipaddr,
		Seq : seq,
		TTL : ttl,
		RTT : rtt.Seconds(),
	}
}

func printPingResult(hostname string, startTime time.Time, transPackets int, statuslist []PingStatus) {
	fmt.Println("\n---", hostname, "ping", "statistics", "---")

	receivePackets := len(statuslist)
	lostPackets := transPackets - receivePackets

	totalTime := time.Now().Sub(startTime)
	
	fmt.Print(transPackets, " packets transmitted, ")
	fmt.Print(receivePackets, " received, ")
	fmt.Printf("%d%% packet loss, time %dms\n", lostPackets*100/transPackets, int(totalTime.Seconds()*1000))

	if receivePackets == 0 {
		return
	}

	min := statuslist[0].RTT
	max := statuslist[0].RTT
	avg := 0.0
	for _, v := range statuslist {
		rtt := v.RTT
		avg += float64(rtt)
		if min > rtt {
			min = rtt
		}
		if max < rtt {
			max = rtt
		}
	}
	avg = avg / float64(receivePackets)

	mdev := 0.0
	for _, v := range statuslist {
		rtt := v.RTT
		mdev += math.Abs(avg - rtt)
	}
	mdev = mdev / float64(receivePackets)

	fmt.Printf("rtt min/avg/max/mdev = %.3f/%.3f/%.3f/%.3f ms\n", min*1000, avg*1000, max*1000, mdev*1000)
}

func main() {

	// Parse Args and Get Host Name
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage : goping <Host Name>")
		os.Exit(1)
	}
	hostname := flag.Arg(0)

	// Get IP address from Host Name
	ip, err := getIPAddr(hostname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create Connection
	connect, err := net.Dial("ip4:1", ip.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer connect.Close()

	// Create ICMP packet
	id := uint16(os.Getpid() & 0xffff)
	sendMessage := &icmp.ICMP{
		Type: icmp.IcmpEchoReq,
		Code: 0,
		Id:   id,
	}
	fmt.Println("PING", hostname, "("+ip.String()+")", "23 bytes of data.")

	// Make ExitNotifier
	exitNotifier := make(chan int, 1)

	// Make channel sigint receive
	sigint := make(chan os.Signal,1)
	signal.Notify(sigint, syscall.SIGINT)

	// Send a ICMP packet
	startTime := time.Now()
	go pinger(connect, sendMessage, exitNotifier, sigint)

	// Receive a ICMP packet
	done := false
	statuslist := make([]PingStatus,0)
	transPackets := 0
	for !done {
		select {
		case transPackets = <-exitNotifier:
			done = true
		default:
			// Read Data
			receiveData := make([]byte, 128)
			connect.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			n, err := connect.Read(receiveData)
			if err != nil {
				continue
			}
			status := getPingStatus(receiveData, n)
			fmt.Println(status)
			statuslist = append(statuslist, status)
		}
	}

	printPingResult(hostname, startTime,transPackets, statuslist)
}

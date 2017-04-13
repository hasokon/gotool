package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
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

func pinger(connection net.Conn, i *icmp.ICMP) {
	seq := uint16(1)

	t := time.NewTicker(1 * time.Second)
	for {
		<-t.C

		timeBinary, err := time.Now().MarshalBinary()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		i.Seq = seq
		i.Data = timeBinary

		_, err = connection.Write(i.Marshal())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		seq++
	}
	t.Stop()
}

func printReceiveData(receiveData []byte, n int) {
	// Parse IP Header
	iphdrLen := uint8(receiveData[0]&0xf) * 4
	iphdr := &ipheader.IPHeader{}
	err := iphdr.Parse(receiveData[:iphdrLen])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//Get IP Addr and Host Name
	dataSize := iphdr.TotalLen
	ttl := iphdr.TTL
	ipaddr := iphdr.SrcAddrString()
	hostname, err := getHostname(ipaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Parse ICMP Message
	receiveData = receiveData[iphdrLen:]
	receiveMessage := &icmp.ICMP{}
	receiveMessage.ParseEchoMessage(receiveData[:n])
	seq := receiveMessage.Seq

	fmt.Printf("%d bytes from %s (%s) : icmp_seq=%d ttl=%d time= ms\n", dataSize, hostname, ipaddr, seq, ttl)
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

	// Send a ICMP packet
	go pinger(connect, sendMessage)

	// Receive a ICMP packet
	for {
		// Read Data
		receiveData := make([]byte, 128)
		n, err := connect.Read(receiveData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		printReceiveData(receiveData, n)
	}
}

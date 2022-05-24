package tools

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type icmp struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func Ping(domain string) (bool, error) {
	originBytes := make([]byte, 2000)
	// 返回一个 ip socket
	//conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
	conn, err := net.DialTimeout("ip4:icmp", domain, 5*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	// 初始化 icmp 报文
	icmpStruct := icmp{8, 0, 0, 0, 0}
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmpStruct)
	//fmt.Println(buffer.Bytes())
	binary.Write(&buffer, binary.BigEndian, originBytes[0:48])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], checkSum(b))
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		return false, err
	}
	return true, nil
}

func checkSum(data []byte) (rt uint16) {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index]) << 8
	}
	rt = uint16(sum) + uint16(sum>>16)
	return ^rt
}

package icmp

import(
	"fmt"
)

const (
	IcmpEchoReq = 8
	IcmpEchoRep = 0
)

type ICMP struct {
	Type uint8
	Code uint8
	CheckSum uint16
	Id uint16
	Seq uint16
	Data []byte
}

func (i *ICMP) String() string {
	return fmt.Sprintf("(%d, %x, %x, %x, %d, DataSize : %d byte)", i.Type, i.Code, i.CheckSum, i.Id, i.Seq, len(i.Data))
}

func (i *ICMP) Marshal() []byte {
	b := make([]byte, 8+len(i.Data))

	// Type
	b[0] = byte(i.Type)

	//Code
	b[1] = byte(i.Code)

	//Check Sum
	b[2] = 0
	b[3] = 0

	//Id Big Engian
	b[4] = byte(i.Id >> 8)
	b[5] = byte(i.Id)

	//Seq Big Endian
	b[6] = byte(i.Seq >> 8)
	b[7] = byte(i.Seq)

	//Data
	copy(b[8:], i.Data)

	//Calculate Check Sum and Set into b
	cs := checksum(b)
	b[2] = byte(cs >> 8)
	b[3] = byte(cs)

	return b
}

func checksum(buf []byte) uint16 {
	len := len(buf)
	sum := uint32(0)
	
	for len > 1 {
		sum += uint32(buf[0]) << 8 | uint32(buf[1])
		buf = buf[2:]
		len -= 2
	}

	if len > 0 {
		sum += uint32(buf[0]) << 8
	}

	sum = (sum & 0xffff) + (sum >> 16)
	sum = (sum & 0xffff) + (sum >> 16)

	return ^(uint16(sum))
}

func (i *ICMP) ParseEchoMessage(b []byte) error {
	i.Type = uint8(b[0])
	i.Code = uint8(b[1])
	i.CheckSum = uint16(b[3]) << 8 + uint16(b[3])
	i.Id = uint16(b[4]) << 8 + uint16(b[5])
	i.Seq = uint16(b[6]) << 8 + uint16(b[7])
	i.Data = make([]byte, len(b[8:]))
	copy(i.Data, b[8:])

	return nil
}

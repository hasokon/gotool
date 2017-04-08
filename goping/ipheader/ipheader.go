package ipheader

import (
	"errors"
)

type IPHeader struct {
	Version uint8 //4bit
	HeaderLen uint8 //4bit
	Service uint8
	TotalLen uint16
	Identification uint16
	Flag uint8 //3bit
	FlagmentOffset uint16 //13bit
	TTL uint8
	Protocol uint8
	CheckSum uint16
	SrcAddr uint32
	DestAddr uint32
	Options uint32
}

func (ih *IPHeader) Parse(b []byte) error {
	len := len(b)
	if len < 20 {
		return errors.New("IP Header Parse Error")
	}

	ih.Version = uint8(b[0] >> 4)
	ih.HeaderLen = uint8(b[0] & 0xf)
	ih.Service = uint8(b[1])
	ih.TotalLen = uint16(b[2] << 8) + uint16(b[3])
	ih.Identification = uint16(b[4] << 8) + uint16(b[5])
	ih.Flag = uint8(b[6] >> 5)
	ih.FlagmentOffset = uint16(b[6] << 8) & 0x1fff + uint16(b[7])
	ih.TTL = uint8(b[8])
	ih.Protocol = uint8(b[9])
	ih.CheckSum = uint16(b[10] << 8) + uint16(b[11])
	ih.SrcAddr = uint32(b[12] << 24) + uint32(b[13] << 16) + uint32(b[14] << 8) + uint32(b[15])
	ih.DestAddr = uint32(b[16] << 24) + uint32(b[17] << 16) + uint32(b[18] << 8) + uint32(b[19])

	if len > 20 {
		ih.Options = uint32(b[20] << 24) + uint32(b[21] << 16) + uint32(b[22] << 8) + uint32(b[23])
	}

	return nil
}

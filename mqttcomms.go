package mqttcomms

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/howeyc/crc16"
)

const (
	pckMagic  = 0x42524144
	magicSize = 4
	seqSize   = 2
	opSize    = 1
	typeSize  = 1
	lenSize   = 2
	crcSize   = 2

	seqOffset  = magicSize
	opOffset   = seqOffset + seqSize
	typeOffset = opOffset + opSize
	lenOffset  = typeOffset + typeSize
	bodyOffset = lenOffset + lenSize
)

const (
	REQACK  = 0x01
	NOACK   = 0x02
	ISACK   = 0x03
	REQRESP = 0x04
	ISRESP  = 0x05
)

const (
	STATUS_OK = 0x0
	STATUS_KO = 0x1
)

type Msg struct {
	Seq    uint16
	Op     uint8
	Method uint8
	Len    uint16
	Body   []byte
	Crc    uint16
}

func GenerateMsg(m *Msg) ([]byte, error) {
	var msg []byte
	buf := new(bytes.Buffer)

	length := len(m.Body)
	if length > 65536 {
		return nil, errors.New("body too big")
	}

	binary.Write(buf, binary.BigEndian, uint32(pckMagic))
	msg = append(msg, buf.Bytes()...)

	buf.Reset()
	binary.Write(buf, binary.BigEndian, uint16(m.Seq))
	msg = append(msg, buf.Bytes()...)
	msg = append(msg, []byte{byte(m.Op)}...)

	msg = append(msg, []byte{byte(m.Method)}...)

	buf.Reset()
	binary.Write(buf, binary.BigEndian, uint16(length))
	msg = append(msg, buf.Bytes()...)

	msg = append(msg, m.Body...)

	crc := crc16.Checksum(msg, crc16.CCITTFalseTable)
	buf.Reset()
	binary.Write(buf, binary.BigEndian, uint16(crc))
	msg = append(msg, buf.Bytes()...)

	return msg, nil
}

func DecodeMsg(msg []byte) (Msg, error) {
	var ret Msg

	if string(msg) == "GOODBYE" {
		return Msg{Method: AD_GOODBYE}, nil
	}

	if len(msg) < 12 {
		return Msg{}, errors.New("message too short")
	}

	if binary.BigEndian.Uint32(msg[0:magicSize]) != uint32(pckMagic) {
		return Msg{}, errors.New("bad magic")
	}

	ret.Seq = binary.BigEndian.Uint16(msg[seqOffset : seqOffset+seqSize])
	ret.Op = uint8(msg[opOffset])
	ret.Method = uint8(msg[typeOffset])
	ret.Len = binary.BigEndian.Uint16(msg[lenOffset : lenOffset+lenSize])

	if len(msg) != bodyOffset+int(ret.Len)+crcSize {
		return Msg{}, errors.New("message has wrong size")
	}

	ret.Body = msg[bodyOffset : bodyOffset+ret.Len]
	ret.Crc = binary.BigEndian.Uint16(msg[bodyOffset+ret.Len : bodyOffset+ret.Len+crcSize])

	computedCrc := crc16.Checksum(msg[0:bodyOffset+ret.Len], crc16.CCITTFalseTable)
	if computedCrc != ret.Crc {
		return Msg{}, errors.New("bad CRC")
	}

	return ret, nil
}

package protocol

import (
	"encoding/binary"
	"io"
)

type SjsonbHeader struct {
	Magic uint32 // 魔数
	Ver uint32 // 版本
	Entype uint16 // 内容类型
	Entofst uint16 // 内容偏移
	Entlen uint32 // 内容长度
	Checksum uint32 // 校验码
}


func SjsonbDecode(rd io.Reader) (*SjsonbHeader, error) {
	var err error
	var header SjsonbHeader

	err = binary.Read(rd, binary.BigEndian, &header)
	if nil != err {
		return nil, err
	}

	return &header, nil
}


// sjsonb协议编码
func SjsonbEncode(ver uint32, entype uint16, cargo []byte) []byte {
	rslt := make([]byte, 20 + len(cargo))

	binary.BigEndian.PutUint32(rslt[0:], 0xE78F8A9D)
	binary.BigEndian.PutUint32(rslt[4:], ver)
	binary.BigEndian.PutUint16(rslt[8:], entype)
	binary.BigEndian.PutUint16(rslt[10:], 20)
	binary.BigEndian.PutUint32(rslt[12:], uint32(len(cargo)))
	binary.BigEndian.PutUint32(rslt[16:], 0)
	copy(rslt[20:], cargo)

	return rslt
}
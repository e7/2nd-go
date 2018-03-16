package protocol

import (
	"bufio"
	"encoding/binary"
	_ "fmt"
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


func SjsonbDecode(rd io.Reader, magic uint32) (*SjsonbHeader, error) {
	var err error
	var header SjsonbHeader
	bufRd := bufio.NewReader(rd)

	// find magic
	for {
		magicBuf, err := bufRd.Peek(4)
		if err != nil {
			return nil, err
		}

		if magic == binary.BigEndian.Uint32(magicBuf) {
			break
		}

		bufRd.Discard(1)
	}

	// decode
	err = binary.Read(bufRd, binary.BigEndian, &header)
	if nil != err {
		return nil, err
	}

	return &header, nil
}


// sjsonb协议编码
func SjsonbEncode(magic, ver uint32, entype uint16, cargo []byte) []byte {
	rslt := make([]byte, 20 + len(cargo))

	binary.BigEndian.PutUint32(rslt[0:], magic)
	binary.BigEndian.PutUint32(rslt[4:], ver)
	binary.BigEndian.PutUint16(rslt[8:], entype)
	binary.BigEndian.PutUint16(rslt[10:], 20)
	binary.BigEndian.PutUint32(rslt[12:], uint32(len(cargo)))
	binary.BigEndian.PutUint32(rslt[16:], 0)
	copy(rslt[20:], cargo)

	return rslt
}

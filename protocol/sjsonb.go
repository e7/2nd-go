package protocol

import (
	"encoding/binary"
	"errors"
	"bytes"
	"bufio"
)

type SjsonbHeader struct {
	Ver uint32 // 版本
	Entype uint16 // 内容类型
	Entofst uint16 // 内容偏移
	Entlen uint32 // 内容长度
	Checksum uint32 // 校验码
}

type Sjsonb struct {
	SjsonbHeader
	Content []byte
}


const (
	ConstIdxProto	  = 4
	ConstIdxFmType    = 4 + 4
	ConstIdxMatterOff = 4 + 4 + 2
	ConstIdxMatterLen = 4 + 4 + 2 + 2
	ConstIdxCheckSum = ConstIdxMatterLen+4

	ConstHdPkgLen = 4 + 4 + 2 + 2 + 4 + 4

	ConstMagic          = 0xE78F8A9D
	ConstVersion uint32 = 1000

	ConstTjson = 0x3
)

func newDummy() *Sjsonb {
	return &Sjsonb{
		SjsonbHeader{
			Ver: 1000,
			Entype:3,
			Entofst:20,
			Entlen:0,
			Checksum:0,
		},
		[]byte{},
	}
}

func New(content []byte) *Sjsonb {
	sjsonb := newDummy()
	sjsonb.Content = content
	sjsonb.Entlen = uint32(len(content))
	return sjsonb
}

func (sjsonb *Sjsonb) Serialize() ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	// magic number
	magic := [4]byte{}
	binary.BigEndian.PutUint32(magic[:], ConstMagic)
	if err := binary.Write(buf, binary.BigEndian, magic[:]); nil != err {
		return nil, err
	}

	// head
	if err := binary.Write(buf, binary.BigEndian, sjsonb.SjsonbHeader); nil != err{
		return nil, err
	}

	// content
	if err := binary.Write(buf, binary.BigEndian, sjsonb.Content); nil != err{
		return nil, err
	}

	//fmt.Println(hex.EncodeToString(buf.Bytes()))
	return buf.Bytes(), nil
}

func SjsonbDecode(rd *bufio.Reader) (*SjsonbHeader, error) {
	var err error
	var header SjsonbHeader

	// find magic
	for {
		magicBuf, err := rd.Peek(4)
		if err != nil {
			return nil, err
		}

		if uint32(ConstMagic) == binary.BigEndian.Uint32(magicBuf) {
			rd.Discard(4) // drop magic
			break
		}

		rd.Discard(1)
	}

	// decode
	err = binary.Read(rd, binary.BigEndian, &header)
	if nil != err {
		return nil, err
	}

	return &header, nil
}

// sjsonb协议编码
func SjsonbEncode(ver uint32, entype uint16, cargo []byte) []byte {
	rslt := make([]byte, 20 + len(cargo))

	binary.BigEndian.PutUint32(rslt[0:], ConstMagic)
	binary.BigEndian.PutUint32(rslt[4:], ver)
	binary.BigEndian.PutUint16(rslt[8:], entype)
	binary.BigEndian.PutUint16(rslt[10:], 20)
	binary.BigEndian.PutUint32(rslt[12:], uint32(len(cargo)))
	binary.BigEndian.PutUint32(rslt[16:], 0)
	copy(rslt[20:], cargo)

	return rslt
}


func Unpack(buffer []byte) (newbuffer []byte, sjsonb *Sjsonb, err error) {
	var i int
	var msgLen int
	var msg16Len uint16
	var msg32Len uint32

	length := len(buffer)

	for i = 0; i < length; i++ {
		if length < i+ConstHdPkgLen {
			break
		}

		if binary.BigEndian.Uint32(buffer[i:i+4]) == ConstMagic {
			msg16Len = binary.BigEndian.Uint16(buffer[i+ConstIdxMatterOff : i+ConstIdxMatterOff+2])
			msg32Len = binary.BigEndian.Uint32(buffer[i+ConstIdxMatterLen : i+ConstIdxMatterLen+4])
			msgLen = int(msg16Len) + int(msg32Len)

			if length < i+msgLen {
				break
			}

			fmt16type := binary.BigEndian.Uint16(buffer[i+ConstIdxFmType : i+ConstIdxFmType+2])
			if fmt16type == ConstTjson {
				content := make([]byte, 0)
				content = append(content, buffer[i+int(msg16Len):i+int(msgLen)]...)
				sjsonb = &Sjsonb{
					SjsonbHeader {
						Ver: binary.BigEndian.Uint32(buffer[i+ConstIdxProto : i+ConstIdxProto+4]),
						Entype: fmt16type,
						Entofst:msg16Len,
						Entlen:msg32Len,
						Checksum:binary.BigEndian.Uint32(buffer[i+ConstIdxCheckSum : i+ConstIdxCheckSum+4]),
						},
					content,
				}
			} else {
				err = errors.New("not support content type")
			}

			i += msgLen
			break
		}
	}

	if i == length {
		newbuffer = make([]byte, 0)
	} else {
		newbuffer = buffer[i:]
	}

	return newbuffer, sjsonb, err
}

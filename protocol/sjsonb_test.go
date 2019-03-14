package protocol

import (
	"testing"
	"bytes"
	"encoding/hex"
	"fmt"
)

func TestSjsonbEncode(t *testing.T) {
	t.Logf("%X", SjsonbEncode(1000, 3, []byte(`{"key":"value"}`)))
}

// 正常情况
func TestSjsonbDecodeNormal(t *testing.T) {
	var err error

	orig := []byte{
		0xe7,0x8f,0x8a,0x9d,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,
	}

	hd, err := SjsonbDecode(bytes.NewReader(orig))
	if nil != err {
		t.Errorf("解码失败:%s", err.Error())
	} else {
		t.Logf("%X", hd.Ver)
		t.Logf("%X", hd.Entype)
		t.Logf("%X", hd.Entofst)
		t.Logf("%X", hd.Entlen)
		t.Logf("%X", hd.Checksum)
		t.Logf("success")
	}
}

// 包含残留数据
func TestSjsonbDecodeWithHeritage(t *testing.T) {
	var err error

	orig := []byte{
		0xab,0xe7,0x8f,0x01,
		0xe7,0x8f,0x8a,0x9d,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,
	}

	hd, err := SjsonbDecode(bytes.NewReader(orig))
	if nil != err {
		t.Errorf("解码失败:%s", err.Error())
	} else {
		t.Logf("%X", hd.Ver)
		t.Logf("%X", hd.Entype)
		t.Logf("%X", hd.Entofst)
		t.Logf("%X", hd.Entlen)
		t.Logf("%X", hd.Checksum)
		t.Logf("success")
	}
}

func TestSjsonbDecodeNoEnoughBytes(t *testing.T) {
	orig2 := []byte{
		0xe7,0x8f,0x8a,0x9d,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,
	} // 不足长度

	hd2, err := SjsonbDecode(bytes.NewReader(orig2))
	if nil != err {
		t.Logf("解码失败:%s\n", err)
	} else {
		t.Logf("%X", hd2.Ver)
		t.Logf("%X", hd2.Entype)
		t.Logf("%X", hd2.Entofst)
		t.Logf("%X", hd2.Entlen)
		t.Logf("%X", hd2.Checksum)
		t.Logf("success")
	}
}

func TestSjsonbUnpack(t *testing.T) {
	teststr := []string{
		"123456E78F8A9D000003E8000300140000000F000000007B226B6579223A2276616C7565227D1234",

		"12E78F8A9D000003E8000300140000000F000000007B226B6579223A2276616C7565227D1234",
		"123456E78F8A9D000003E8000300140000000F000000007B226B6579223A2276616C7565227D1234",
	}
	for _, str := range teststr {
		bytes1, _ := hex.DecodeString(str)

		leave, sjsonb, err := Unpack(bytes1)
		if nil != err || 2 != len(leave) {
			t.Errorf("unpack failed")
		}

		if len(sjsonb.Content) != int(sjsonb.Entlen) {
			t.Errorf("sjonb err")
		}
	}

	errstr := []string {
		"E88F8A9D000003E8000300140000000F000000007B226B6579223A2276616C7565227D1234",
		"123456",
	}

	for _, str := range errstr {
		bytes1, _ := hex.DecodeString(str)

		leave, sjsonb, err := Unpack(bytes1)
		if nil != sjsonb || nil != err || len(leave) == 0 {
			t.Errorf("unpack failed")
		}
		fmt.Println(hex.EncodeToString(leave))
	}
}

func TestSjsonbSerialize(t *testing.T) {
	sjonb := New(nil)
	sjbuf, err := sjonb.Serialize()
	if nil != err {
		t.Errorf("serialize sjsonb failed!")
	}

	leave, sjsonb, err := Unpack(sjbuf)
	if nil == sjsonb || nil != err || len(leave) != 0 {
		t.Errorf("unpack failed")
	}

	fmt.Println(sjsonb)
}

func TestSjsonbSerialize1(t *testing.T) {
	sjonb := New([]byte{1,2,3})
	sjbuf, err := sjonb.Serialize()
	if nil != err {
		t.Errorf("serialize sjsonb failed!")
	}

	leave, sjsonb, err := Unpack(sjbuf)
	if nil == sjsonb || nil != err || len(leave) != 0 {
		t.Errorf("unpack failed")
	}

	fmt.Println(sjsonb)
}
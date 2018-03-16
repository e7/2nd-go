package protocol

import (
	"testing"
	"bytes"
	"unsafe"
)

func TestSjsonbEncode(t *testing.T) {
	t.Logf("%X", SjsonbEncode(
		0xE78F8A9D,1000, 3, []byte(`{"key":"value"}`),
	))
}

func TestSjsonbDecode(t *testing.T) {
	var err error

	t.Logf("sizeof SjsonbHeader:%d", unsafe.Sizeof(SjsonbHeader{}))

	orig := []byte{
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,
	}

	hd, err := SjsonbDecode(bytes.NewReader(orig), 0x12345678)
	if nil != err {
		t.Errorf("解码失败:%s", err.Error())
	} else {
		t.Logf("%X", hd.Magic)
		t.Logf("%X", hd.Ver)
		t.Logf("%X", hd.Entype)
		t.Logf("%X", hd.Entofst)
		t.Logf("%X", hd.Entlen)
		t.Logf("%X", hd.Checksum)
		t.Logf("success")
	}

	t.Log("--------------------------------------------")

	orig2 := []byte{
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,
	} // 不足长度

	hd2, err := SjsonbDecode(bytes.NewReader(orig2), 0x12345678)
	if nil != err {
		t.Logf("解码失败:%s\n", err)
	} else {
		t.Errorf("%X", hd2.Magic)
		t.Errorf("%X", hd2.Ver)
		t.Errorf("%X", hd2.Entype)
		t.Errorf("%X", hd2.Entofst)
		t.Errorf("%X", hd2.Entlen)
		t.Errorf("%X", hd2.Checksum)
		t.Errorf("success")
	}

	t.Log("--------------------------------------------")

	orig3 := []byte{
		0xe7,0x8f,0x8a,0x9d,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,0xbc,0xde,0xf0,
		0x12,0x34,0x56,0x78,0x9a,
	} // 长度超出

	hd3, err := SjsonbDecode(bytes.NewReader(orig3), 0xe78f8a9d)
	if nil != err {
		t.Errorf("解码失败:%s", err)
	} else {
		t.Logf("%X", hd3.Magic)
		t.Logf("%X", hd3.Ver)
		t.Logf("%X", hd3.Entype)
		t.Logf("%X", hd3.Entofst)
		t.Logf("%X", hd3.Entlen)
		t.Logf("%X", hd3.Checksum)
		t.Logf("success")
	}
}

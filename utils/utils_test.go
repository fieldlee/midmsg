package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	b := StringToBytes("ent2015")
	fmt.Println(b)

	s := BytesToString(b)
	fmt.Println(s)
}

func TestInt16ToBytes(t *testing.T) {
	b:=Int16ToBytes(1000)
	fmt.Println(b)
	v := BytesToInt16(b)
	fmt.Println(v)
}

func TestUint8ToBytes(t *testing.T) {
	b := Uint8ToBytes(1)
	fmt.Println(b)
	v := BytesToUInt8(b)
	fmt.Println(v)
}

func TestBytesJoin(t *testing.T) {
	b1 := Int16ToBytes(1000)
	fmt.Println(b1)
	b2 := Uint8ToBytes(1)
	fmt.Println(b2)
	b3 := Uint32ToBytes(3000)
	fmt.Println(b3)
	b := BytesJoin(b1,b2,b3)
	fmt.Println(b)


}
func TestPrintCh(t *testing.T) {
	PrintCh()
	b := []byte("hello我爱你中国！")
	BytesTOCh(b)

	list := strings.FieldsFunc("你好, 我是 李德鹏, 你是 哪位", func(r rune) bool {
		if r == '是'{
			return true
		}
		return false
	})
	fmt.Print(list)
}

func TestEncryptDecrypt(t *testing.T) {
	src := []byte("hello world!")
	key := []byte("meimeigujiagujia")
	encryptByte ,err := EncryptAes(src,key)
	if err != nil {
		fmt.Println(err)
	}
	dsrc,err := DecryptAes(encryptByte,key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dsrc))
}

func TestZipBytes(t *testing.T) {
	src := []byte("hello worldasdfasdfasdfasdfasdfasdfadsfasdfasdfasdfasdfasdfasdfasdfasdfasdfsdfasdfasdfasdfasdfasdfasdf!")
	fmt.Println(len(src))
	zipsrc,err := ZipByte(src)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(zipsrc))
	unzipbyte,err  := UnzipByte(zipsrc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(unzipbyte),len(unzipbyte))
}
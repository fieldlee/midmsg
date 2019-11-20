package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
)

func ClearBytes(origin []byte)[]byte{
	x := len(origin)
	tran := make([]byte,x)
	for i,b := range origin{
		t := fmt.Sprintf("%v",b)
		if t != "0" {
			tran[i] = b
		}
	}
	return tran
}

func BytesToString(b []byte)string{
	nb := make([]byte,0)
	for _,t := range b {
		x := fmt.Sprintf("%v",t)
		if x != "0"{
			nb = append(nb,t)
		}
	}
	return string(nb)
}

//字节转换成整形
func BytesToInt16(b []byte) int16 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}
// unsigned char -->  C.uchar -->  uint8
func BytesToUInt8(b []byte) uint8 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x uint8
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}

func BytesToInt32(b []byte) int32 {
	//b := ClearBytes(by)
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int32(x)
}

func UInt16ToBytes(i uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

func UnzipBytes(zip []byte)[]byte{
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	defer w.Close()
	w.Write(zip[:])
	w.Flush()
	r, err := gzip.NewReader(&b)
	if err != nil {
		return zip
	}
	defer r.Close()
	undatas, err := ioutil.ReadAll(r)
	if err != nil {
		return zip
	}
	return undatas
}

func Goid() int {
	defer func()  {
		if err := recover(); err != nil {
			fmt.Println("panic recover:panic info:%v", err)        }
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
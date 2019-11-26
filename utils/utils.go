package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc/peer"
	"io/ioutil"
	"midmsg/model"
	"net"
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
//
//func Goid() int {
//	defer func()  {
//		if err := recover(); err != nil {
//			fmt.Println("panic recover:panic info:%v", err)        }
//	}()
//
//	var buf [64]byte
//	n := runtime.Stack(buf[:], false)
//	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
//	id, err := strconv.Atoi(idField)
//	if err != nil {
//		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
//	}
//	return id
//}

func Decrypt(b []byte,encrptType model.ENCRPTION_TYPE)[]byte{
	if encrptType == model.Encryption_AES{

	}
	if encrptType == model.Encryption_Des{

	}
	if encrptType == model.Encryption_RSA{

	}
	return b
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func padding(src []byte,blocksize int) []byte {
	padnum:=blocksize-len(src)%blocksize
	pad:=bytes.Repeat([]byte{byte(padnum)},padnum)
	return append(src,pad...)
}

func unpadding(src []byte) []byte {
	n:=len(src)
	unpadnum:=int(src[n-1])
	return src[:n-unpadnum]
}

func Encrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	src=padding(src,block.BlockSize())
	blockmode:=cipher.NewCBCEncrypter(block,key[:block.BlockSize()])
	blockmode.CryptBlocks(src,src)
	return src
}

func Decrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	blockmode:=cipher.NewCBCDecrypter(block,key[:block.BlockSize()])
	blockmode.CryptBlocks(src,src)
	src=unpadding(src)
	return src
}


func GetClietIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("getClinetIP, invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("getClientIP, peer.Addr is nil")
	}

	if strings.Contains(pr.Addr.String(),":"){
		return strings.Split(pr.Addr.String(),":")[0],nil
	}

	return pr.Addr.String(), nil
}
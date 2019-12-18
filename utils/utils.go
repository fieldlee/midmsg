package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"google.golang.org/grpc/peer"
	"io/ioutil"
	"midmsg/model"
	"net"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
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

func JoinHeadAndBody(info model.HeadInfo,in []byte)[]byte{
	tagByte := StringToBytes(info.Tag)
	versionByte := Int16ToBytes(info.Version)
	clientTypeByte := Int16ToBytes(info.ClientType)
	headLengthByte := Int16ToBytes(info.HeadLength)
	CompressWayBYte := Uint8ToBytes(info.CompressWay)
	EncryptionBYte := Uint8ToBytes(info.Encryption)
	SigBYte := Uint8ToBytes(info.Sig)
	FormatBYte := Uint8ToBytes(info.Format)
	NetFlagBYte := Uint8ToBytes(info.NetFlag)
	Back1BYte := Uint8ToBytes(info.Back1)
	BufSizeBYte := Int32ToBytes(info.BufSize)
	UncompressedSizeByte := Int32ToBytes(info.UncompressedSize)
	Back2Byte := Int32ToBytes(info.Back2)
	return BytesJoin(tagByte,versionByte,clientTypeByte,headLengthByte,CompressWayBYte,EncryptionBYte,SigBYte,
		FormatBYte,NetFlagBYte,Back1BYte,BufSizeBYte,UncompressedSizeByte,Back2Byte,in)
}

func BytesJoin(blist ...[]byte)[]byte{
	bytesinfo := make([]byte,0)
	for _,b := range blist  {
		bytesinfo = append(bytesinfo,b...)
	}
	return bytesinfo
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

func StringToBytes(s string)[]byte{
	b := []byte(s)
	return b[:8]
}


func Int16ToBytes(n int16)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:2]
}

func Uint8ToBytes(n uint8)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:1]
}

func Uint32ToBytes(n uint32)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:4]
}

func Int32ToBytes(n int32)[]byte{
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.LittleEndian, n)
	if err != nil {
		return nil
	}
	return buffer.Bytes()[:4]
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

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func Decrypt(b []byte,encrptType model.ENCRPTION_TYPE)[]byte{
	if encrptType == model.Encryption_AES{

	}
	if encrptType == model.Encryption_Des{

	}
	if encrptType == model.Encryption_RSA{

	}
	return b
}

func DecryptAes(crypted, key []byte) ([]byte, error) {
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

func EncryptAes(origData, key []byte) ([]byte, error) {
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
	block,_:= des.NewTripleDESCipher(key)
	src= padding(src,block.BlockSize())
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

func EncryptRsa(originalData,key []byte)([]byte,error){
	pubKey, _ := x509.ParsePKIXPublicKey(key)
	encryptedData,err:=rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), originalData)
	return encryptedData,err
}

//（2）解密：对采用sha1算法加密后转base64格式的数据进行解密（私钥PKCS1格式）
func DecryptRsa(encryptedData,key []byte)([]byte,error){
	prvKey,_:=x509.ParsePKCS1PrivateKey(key)
	originalData,err:=rsa.DecryptPKCS1v15(rand.Reader,prvKey,encryptedData)
	return originalData,err
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

func PrintCh(){
	ch := "hello我爱你中国！"
	for i,v := range ch {
		// v is rune 类型
		fmt.Printf("(%d %X)",i,v)
	}
	fmt.Println()
}
func BytesTOCh(b []byte){
	for len(b)>0{
		ch,s := utf8.DecodeRune(b)
		fmt.Printf("%c",ch)
		b = b[s:]
	}
}

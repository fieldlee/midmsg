package handle

import (

	"context"
	"encoding/json"
	"fmt"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
)

func init()  {
	JobQueue = make(chan HandleBody, utils.MaxQueue)
}
type MsgHandle struct {}

func (m *MsgHandle)Sync(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo, error) {

	out := make(chan *pb.NetRspInfo)
	err := make(chan error)
	//// 发送body到队列
	handleBody := HandleBody{
		M_Body:in.M_Body,
		Out: out,
		Err:err,
	}

	JobQueue <- handleBody

	for {
		select {
		case netrep := <-out:
			return netrep,nil
		case errinfo := <-err:
			return nil,errinfo
		}
	}
}

func (m *MsgHandle)Async(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo, error) {
	out := new(pb.NetRspInfo)

	inByte := in.M_Body
	// 解析body bytes 头 32个字节
	err := AnzalyBodyHead(inByte)

	if err != nil {
		return out , err
	}

	return out , nil
}


func AnzalyBodyHead(inbody []byte) error {
	bodyHead := inbody[:32]
	//包头标示  8
	m_tag := bodyHead[:8]
	bodyHead = bodyHead[8:]

	tag := utils.BytesToString(m_tag)
	fmt.Println("tag:",tag)
	//数据版本  2
	m_lDateVersion := bodyHead[:2]
	bodyHead = bodyHead[2:]
	version := utils.BytesToInt16(m_lDateVersion)
	fmt.Println("version:",version)
	//客户端类型 2
	m_sClientType := bodyHead[:2]
	bodyHead = bodyHead[2:]
	clientType := utils.BytesToInt16(m_sClientType)
	fmt.Println("clientType:",clientType)
	/////check client type
	if clientType >= int16(model.ClientTypeMax) {
		return model.ErrClientType
	}
	//包头长度   2
	m_sHeadLength := bodyHead[:2]
	bodyHead = bodyHead[2:]
	headLength := utils.BytesToInt16(m_sHeadLength)
	fmt.Println("headLength:",headLength)
	////check head length
	if headLength != 32 {
		return model.ErrHeaderLength
	}
	//压缩方式   1
	m_cCompressionWay := bodyHead[:1]
	bodyHead = bodyHead[1:]
	compressionWay := utils.BytesToUInt8(m_cCompressionWay)
	fmt.Println("compressionWay:",compressionWay)
	if compressionWay >= uint8(model.CompressionWayMax) {
		return model.ErrCompressionType
	}
	//加密方式   1
	m_cEncryption := bodyHead[:1]
	bodyHead = bodyHead[1:]
	encryption := utils.BytesToUInt8(m_cEncryption)
	fmt.Println("encryption:",encryption)
	if encryption >= uint8(model.Encryption_Max) {
		return model.ErrEncrptyType
	}
	//协议标识   1
	m_cSig := bodyHead[:1]
	bodyHead = bodyHead[1:]
	sig := utils.BytesToUInt8(m_cSig)
	fmt.Println("sig:",sig)
	//数据流格式  1
	m_cdataFormat := bodyHead[:1]
	bodyHead = bodyHead[1:]
	format := utils.BytesToUInt8(m_cdataFormat)
	fmt.Println("format:",format)
	//网络标记   1
	m_cNetFlag := bodyHead[:1]
	bodyHead = bodyHead[1:]
	flag := utils.BytesToUInt8(m_cNetFlag)
	fmt.Println("flag:",flag)
	//占位符     1
	m_cBack1 := bodyHead[:1]
	bodyHead = bodyHead[1:]
	fmt.Println("占位符1:",utils.BytesToUInt8(m_cBack1))
	//数据长度   4
	m_lBufSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	bufSize := utils.BytesToInt32(m_lBufSize)
	fmt.Println("bufSize:",bufSize)
	///// 校验数据长度
	if int32(len(inbody)-32) != bufSize {
		return model.ErrCompressedLength
	}

	//压缩前长度 4
	m_lUncompressedSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	uncompressiondSize := utils.BytesToInt32(m_lUncompressedSize)
	fmt.Println("compressiondSize:",uncompressiondSize)
	if bufSize > uncompressiondSize {
		return model.ErrCompresseduncompressedLength
	}
	//////如果压缩，解压然后比较长度

	if compressionWay == uint8(model.Compression_zip){
		unzipBytes:=utils.UnzipBytes(inbody[32:])
		if uncompressiondSize != int32(len(unzipBytes)){
			return model.ErrUNCompressedLength
		}
	}

	//预留位     4
	m_lBack2 := bodyHead[:]
	fmt.Println("占位符2:",utils.BytesToInt32(m_lBack2))
	return nil
}

func AnzalyBody(inbody []byte) (*model.GJ_Net_Pack,error) {
	body := inbody[31:]
	netPack := model.GJ_Net_Pack{}
	//netPack := make(map[interface{}]interface{})
	err := json.Unmarshal(body,&netPack)
	if err != nil {
		return nil,err
	}

	fmt.Println("netPack:",netPack)

	return &netPack,nil
}
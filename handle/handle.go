package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
)

type MsgHandle struct {
	
}

func (m *MsgHandle)Sync(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo, error) {
	out := new(pb.NetRspInfo)

	inByte := in.M_Body
	// 解析body bytes 头 32个字节
	err := AnzalyBodyHead(inByte)

	if err != nil {
		return out , err
	}

	return out ,nil
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
	//数据版本  2
	m_lDateVersion := bodyHead[:2]
	bodyHead = bodyHead[2:]
	version := utils.BytesToInt(m_lDateVersion)
	//客户端类型 2
	m_sClientType := bodyHead[:2]
	bodyHead = bodyHead[2:]
	clientType := utils.BytesToInt(m_sClientType)
	//包头长度   2
	m_sHeadLength := bodyHead[:2]
	bodyHead = bodyHead[2:]
	headLength := utils.BytesToInt(m_sHeadLength)
	//压缩方式   1
	m_cCompressionWay := bodyHead[:1]
	bodyHead = bodyHead[1:]
	compressionWay := utils.BytesToString(m_cCompressionWay)
	//加密方式   1
	m_cEncryption := bodyHead[:1]
	bodyHead = bodyHead[1:]
	encryption := utils.BytesToString(m_cEncryption)
	//协议标识   1
	m_cSig := bodyHead[:1]
	bodyHead = bodyHead[1:]
	sig := utils.BytesToString(m_cSig)
	//数据流格式  1
	m_cdataFormat := bodyHead[:1]
	bodyHead = bodyHead[1:]
	format := utils.BytesToString(m_cdataFormat)
	//网络标记   1
	m_cNetFlag := bodyHead[:1]
	bodyHead = bodyHead[1:]
	flag := utils.BytesToString(m_cNetFlag)
	//占位符     1
	m_cBack1 := bodyHead[:1]
	bodyHead = bodyHead[1:]
	fmt.Println("占位符1:",string(m_cBack1))
	//数据长度   4
	m_lBufSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	bufSize := utils.BytesToInt(m_lBufSize)
	//压缩前长度 4
	m_lUncompressedSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	compressiondSize := utils.BytesToInt(m_lUncompressedSize)
	//预留位     4
	m_lBack2 := bodyHead[:4]
	fmt.Println("占位符2:",string(m_lBack2))
	return nil
}

func AnzalyBody(inbody []byte) error {
	body := inbody[32:]
	netPack := model.GJ_Net_Pack{}
	err := json.Unmarshal(body,&netPack)
	if err != nil {
		return err
	}

	return nil
}
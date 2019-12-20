package handle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/pborman/uuid"
	"midmsg/call"
	"midmsg/log"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
	"strings"
	"sync"
	"time"
)

var SqlClient  *utils.SqlCliet
var sqlerr error
var SeqIP = map[string]string{}
var SubScribeDetail = map[string][]string{}
func init()  {
	SqlClient,sqlerr = utils.InitSql()
	if sqlerr != nil {
		log.Fatal(sqlerr)
	}
	loadSeqIP()
	loadSubScribe()
}

func loadSeqIP(){
	SeqIP,sqlerr = SqlClient.GetAllFunc()
	if sqlerr != nil {
		log.Fatal(sqlerr)
	}
}

func loadSubScribe(){
	SubScribeDetail,sqlerr = SqlClient.GetAllSubScribe()
	if sqlerr != nil {
		log.Fatal(sqlerr)
	}
}

type MsgHandle struct {}

func (m *MsgHandle)Sync(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo, error) {
	ipaddr,err := utils.GetClietIP(ctx)
	if err != nil {
		return nil,err
	}

	out := make(chan *pb.NetRspInfo)
	//// 发送body到队列
	handleBody := HandleBody{
		ClientIp:ipaddr,
		MBody:in.M_Body,
		Type: model.CALL_CLIENT_SYNC,
		Out: out,
	}
	///// 异步answer
	if call.CheckAsyncAnswer(in.Uuid) {
		handleBody = HandleBody{
			Sequence:in.Uuid,
			ClientIp:ipaddr,
			MBody:in.M_Body,
			Type: model.CALL_CLIENT_ANSWER,
			Out: out,
		}
	}

	go func(handleBody HandleBody) {
		JobQueue <- handleBody
	}(handleBody)

	for {
		select {
		case _ = <-JobDone:
			return <-handleBody.Out,nil
		}
	}
}

func (m *MsgHandle)Async(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo, error) {
	ipaddr,err := utils.GetClietIP(ctx)
	if err != nil {
		return nil,err
	}

	out := make(chan *pb.NetRspInfo)
	//// 发送body到队列
	handleBody := HandleBody{
		ClientIp:ipaddr,
		MBody:in.M_Body,
		Type: model.CALL_CLIENT_ASYNC,
		Out: out,
	}

	go func(handleBody HandleBody) {
		JobQueue <- handleBody
	}(handleBody)

	for {
		select {
		case _ = <-JobDone:
			return <-handleBody.Out,nil
		}
	}
}

func (m *MsgHandle)Broadcast(ctx context.Context, in *pb.NetReqInfo) (*pb.NetRspInfo,error){
	ipaddr,err := utils.GetClietIP(ctx)
	if err != nil {
		return nil,err
	}

	out := make(chan *pb.NetRspInfo)
	//// 发送body到队列
	handleBody := HandleBody{
		Service:in.Service,
		ClientIp:ipaddr,
		MBody:in.M_Body,
		Type: model.CALL_CLIENT_PUBLISH,
		Out: out,
	}

	go func(handleBody HandleBody) {
		JobQueue <- handleBody
	}(handleBody)

	for {
		select {
		case _ = <-JobDone:
			return <-handleBody.Out,nil
		}
	}
}

func (m *MsgHandle)Register(ctx context.Context, in *pb.RegisterInfo)(*pb.RegisterReturnInfo,error){
	funid,err := SqlClient.GetFunc(in.Sequence)
	if err != nil{
		return nil,err
	}
	if funid != ""{
		return &pb.RegisterReturnInfo{
			Success:true,
		},errors.New(fmt.Sprintf("the %s function had registered",in.Sequence))
	}
	/////保存funcid
	if _,err := SqlClient.GetFunc(in.Sequence);err != nil && err != sql.ErrNoRows {
		return nil,err
	}else{
		err = SqlClient.InsertFunc(in.Sequence)
		if err != nil{
			return nil,err
		}
	}


	clientIP,err := utils.GetClietIP(ctx)
	if err != nil{
		return nil,err
	}

	////保存

	if ip,err := SqlClient.GetFuncListByIP(in.Sequence,fmt.Sprintf("%s:%d",clientIP,utils.ClientPort));err != nil && err!= sql.ErrNoRows{
		return nil,err
	}else{
		if ip == ""{
			err = SqlClient.InsertFuncList(in.Sequence,fmt.Sprintf("%s:%d",clientIP,utils.ClientPort))
			if err != nil{
				return nil,err
			}
		}else{
			return &pb.RegisterReturnInfo{
				Success:true,
			},errors.New(fmt.Sprintf("the %s function had registered",in.Sequence))
		}
	}

	loadSeqIP()
	log.Error(SeqIP)
	return &pb.RegisterReturnInfo{
		Success:true,
	}, nil
}

func (m *MsgHandle)Publish(ctx context.Context, in *pb.PublishInfo)(*pb.PublishReturnInfo,error){
	svcid,err := SqlClient.GetSvc(in.Service)
	if err != nil{
		return nil,err
	}
	if svcid != ""{
		return &pb.PublishReturnInfo{
			Success:true,
		},errors.New(fmt.Sprintf("the %s service had registered",in.Service))
	}

	if _,err := SqlClient.GetSvc(in.Service);err != nil && err != sql.ErrNoRows{
		return nil,err
	}else{
		err = SqlClient.InsertSvc(in.Service)
		if err != nil {
			return nil,err
		}
	}

	return &pb.PublishReturnInfo{
		Success:true,
	},nil
}

func (m *MsgHandle)Subscribe(ctx context.Context, in *pb.SubscribeInfo)(*pb.SubscribeReturnInfo,error) {
	clientIP,err := utils.GetClietIP(ctx)
	if err != nil{
		return nil,err
	}
	ip,err := SqlClient.GetSubScribeByIP(in.Service,fmt.Sprintf("%s:%d",clientIP,utils.ClientPort))
	if err != nil{
		return nil,err
	}
	if ip != ""{
		return &pb.SubscribeReturnInfo{
			Success:true,
		},errors.New(fmt.Sprintf("the %s service and %s ip had registered",in.Service,fmt.Sprintf("%s:%d",clientIP,utils.ClientPort)))
	}
	err = SqlClient.InsertSubScribe(in.Service,fmt.Sprintf("%s:%d",clientIP,utils.ClientPort))
	if err != nil {
		return nil,err
	}

	loadSubScribe()
	log.Error(SubScribeDetail)
	return &pb.SubscribeReturnInfo{
		Success:true,
	},nil
}

func CheckHaveHead(inbody []byte) bool{
	body := inbody[32:]
	netPack := pb.GJ_Net_Pack{}
	errAll := proto.Unmarshal(inbody,&netPack)
	errBak :=  proto.Unmarshal(body,&netPack)
	if errAll==nil && errBak != nil {
		return false
	}
	return true
}

func ModifyBody(inbody []byte, info *model.HeadInfo)[]byte{
	bodyContent := inbody[32:]
	//////如果是压缩，解压缩后再比较
	if model.COMPRESS_TYPE(info.CompressWay) == model.Compression_zip{
		bodyContent,_ = utils.UnzipByte(bodyContent)
	}
	//////解析加密和解密
	switch model.ENCRPTION_TYPE(info.Encryption) {
	case model.Encryption_Des:
		return utils.Decrypt3DES(bodyContent,[]byte(model.PassPass24))
	case model.Encryption_AES:
		aesByte,_ := utils.DecryptAes(bodyContent,[]byte(model.PassPass16))
		return aesByte
	case model.Encryption_RSA:
		prikey := utils.BytesToPrivateKey(model.PriKeyByte)
		rsaByte := utils.DecryptWithPrivateKey(bodyContent,prikey)
		return rsaByte
	default:
		return bodyContent
	}
}

func AnzalyBodyHead(inbody []byte) (*model.HeadInfo,error) {
	headinfo := model.HeadInfo{}
	if len(inbody) < 32 { ////校验包的长度，必须大于32
		return nil,fmt.Errorf("the package bytes length too short.")
	}
	bodyHead := inbody[:32]
	//包头标示  8
	m_tag := bodyHead[:8]
	bodyHead = bodyHead[8:]

	tag := utils.BytesToString(m_tag)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"tag:",tag)
	headinfo.Tag = tag

	//数据版本  2
	m_lDateVersion := bodyHead[:2]
	bodyHead = bodyHead[2:]
	version := utils.BytesToInt16(m_lDateVersion)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"version:",version)
	headinfo.Version = version
	//客户端类型 2
	m_sClientType := bodyHead[:2]
	bodyHead = bodyHead[2:]
	clientType := utils.BytesToInt16(m_sClientType)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"clientType:",clientType)
	/////check client type
	if clientType >= int16(model.ClientTypeMax) {
		return nil,model.ErrClientType
	}
	headinfo.ClientType = clientType
	//包头长度   2
	m_sHeadLength := bodyHead[:2]
	bodyHead = bodyHead[2:]
	headLength := utils.BytesToInt16(m_sHeadLength)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"headLength:",headLength)
	//check head length
	if headLength != 32 {
		return nil,model.ErrHeaderLength
	}
	headinfo.HeadLength = headLength
	//压缩方式   1
	m_cCompressionWay := bodyHead[:1]
	bodyHead = bodyHead[1:]
	compressionWay := utils.BytesToUInt8(m_cCompressionWay)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"compressionWay:",compressionWay)

	if compressionWay >= uint8(model.CompressionWayMax) {
		return nil,model.ErrCompressionType
	}
	headinfo.CompressWay = compressionWay


	//加密方式   1
	m_cEncryption := bodyHead[:1]
	bodyHead = bodyHead[1:]
	encryption := utils.BytesToUInt8(m_cEncryption)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"encryption:",encryption)
	if encryption >= uint8(model.Encryption_Max) {
		return nil,model.ErrEncrptyType
	}
	headinfo.Encryption = encryption
	//协议标识   1
	m_cSig := bodyHead[:1]
	bodyHead = bodyHead[1:]
	sig := utils.BytesToUInt8(m_cSig)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"sig:",sig)
	headinfo.Sig = sig
	//数据流格式  1
	m_cdataFormat := bodyHead[:1]
	bodyHead = bodyHead[1:]
	format := utils.BytesToUInt8(m_cdataFormat)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"format:",format)
	headinfo.Format = format
	//网络标记   1
	m_cNetFlag := bodyHead[:1]
	bodyHead = bodyHead[1:]
	flag := utils.BytesToUInt8(m_cNetFlag)

	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"flag:",flag)
	headinfo.NetFlag = flag
	//占位符     1
	m_cBack1 := bodyHead[:1]
	bodyHead = bodyHead[1:]
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"占位符1:",utils.BytesToUInt8(m_cBack1))
	headinfo.Back1 = utils.BytesToUInt8(m_cBack1)
	//数据长度   4
	m_lBufSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	bufSize := utils.BytesToInt32(m_lBufSize)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"bufSize:",bufSize)
	///// 校验数据长度
	if int32(len(inbody)-32) != bufSize {
		return nil,model.ErrCompressedLength
	}
	headinfo.BufSize = bufSize
	//压缩前长度 4
	m_lUncompressedSize := bodyHead[:4]
	bodyHead = bodyHead[4:]
	uncompressiondSize := utils.BytesToInt32(m_lUncompressedSize)
	//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBodyHead"},"uncompressiondSize:",uncompressiondSize)
	//////如果压缩，解压然后比较长度
	if model.COMPRESS_TYPE(compressionWay) == model.Compression_zip{
		unzipBytes,_:=utils.UnzipByte(inbody[32:])
		if uncompressiondSize != int32(len(unzipBytes)){
			return nil,model.ErrUNCompressedLength
		}
	}else{
		/////不压缩的数据长度相同
		if bufSize != uncompressiondSize{
			return nil,model.ErrUNCompressedLength
		}
	}

	headinfo.UncompressedSize = uncompressiondSize
	//预留位     4
	m_lBack2 := bodyHead[:]

	headinfo.Back2 = utils.BytesToInt32(m_lBack2)
	return &headinfo,nil
}

func AnzalyBody(inbody []byte,uuid string,syncType model.CALL_CLIENT_TYPE,clientIP string) (*pb.NetRspInfo,error) {
	body := inbody
	netPack := pb.GJ_Net_Pack{}
	err :=  proto.Unmarshal(body,&netPack)
	if err != nil {
		return &pb.NetRspInfo{
			M_Err:[]byte(err.Error()),
		},nil
	}

	collectResult := map[uint32]*pb.SendResultInfo{}

	singleResult := make(chan pb.SendResultInfo,len(netPack.M_Net_Pack))

	for key,_ := range netPack.M_Net_Pack {
		if syncType == model.CALL_CLIENT_ANSWER { /////异步回答
			go AsyncAnswer(key ,netPack.M_Net_Pack[key],uuid,syncType,clientIP,singleResult)
		}else{  /////同步请求
			go CheckAndSend(key ,netPack.M_Net_Pack[key],uuid,syncType,clientIP,singleResult)
		}
	}
	defer close(singleResult)
	/////读取返回值
	for i := 0;i<len(netPack.M_Net_Pack);i++{
		tmpResult := <- singleResult
		collectResult[tmpResult.Key] = &tmpResult
	}
	return &pb.NetRspInfo{
		M_Net_Rsp:collectResult,
	},nil
}

func CheckAndSend(key uint32,netpack *pb.Net_Pack,suuid string,syncType model.CALL_CLIENT_TYPE,clientIP string,result chan pb.SendResultInfo){
	tSendResult := pb.SendResultInfo{
		Key:key,
		Result:nil,
		CheckErr:nil,
	}
	///check ASK_TYPE
	if netpack.M_MsgBody.MLAsktype > uint64(model.ETN_SERVER_SUBSRCTIBE_MSG) {
		tSendResult.CheckErr = []byte(model.ErrAskType.Error())
		result <- tSendResult
		return
	}
	//// check Msg——type
	if netpack.M_MsgBody.MCMsgType >= int32(model.MSG_TYPEMAX){
		tSendResult.CheckErr = []byte(model.ErrMsgType.Error())
		result <- tSendResult
		return
	}
	//// check Send count
	if netpack.M_MsgBody.MSSendCount < 1 {
		tSendResult.CheckErr = []byte(model.ErrSendCount.Error())
		result <- tSendResult
		return
	}
	// 失败了，是否要丢弃
	isDiscard := false
	if netpack.M_MsgBody.MIDiscard == 0 {
		isDiscard = true
	}
	/******
	/// select ASK_TYPE
	sevices := utils.GetServiceByKey(fmt.Sprintf("%d",netpack.M_MsgBody.MLAsktype))
	address := sevices["address"].(string)
	port 	:= fmt.Sprintf("%d",sevices["port"])
	service := sevices["service"].(string)
	 */
	tempIP := SeqIP[fmt.Sprintf("%d",netpack.M_MsgBody.MLAsktype)]

	if tempIP == "" {
		tSendResult.CheckErr = []byte(errors.New(fmt.Sprintf("the %d sequence not found address",netpack.M_MsgBody.MLAsktype)).Error())
		result <- tSendResult
		return
	}
	address := strings.Split(tempIP,":")[0]
	port := strings.Split(tempIP,":")[1]
	service := ""

	///// 超时时间
	timeout :=  time.Second * time.Duration(netpack.M_MsgBody.MLExpireTime)

	sendInfo := model.CallInfo{
		Sequence:uuid.New(),
		ClientIP:clientIP,
		Address:address,
		Port:port,
		Service:service,
		MsgBody:netpack,
		Timeout:timeout,
		IsDiscard:isDiscard,
		AskSequence:netpack.M_MsgBody.MLAskSequence,
		SendTimeApp:netpack.M_MsgBody.MISendTimeApp,
		MsgType :netpack.M_MsgBody.MCMsgType,
		MsgAckType :netpack.M_MsgBody.MCMsgAckType,
		SyncType:syncType,
	}

	rspinfo,err := call.CallClient(sendInfo)

	if err != nil {
		tSendResult.CheckErr = []byte(err.Error())
	}

	tSendResult.CheckErr = nil
	tSendResult.Result = rspinfo.M_Net_Rsp

	result <- tSendResult
	return
}

func PublishBody(inbody []byte,service,clientIP string) (*pb.NetRspInfo,error) {
	body := inbody
	netPack := pb.GJ_Net_Pack{}
	err :=  proto.Unmarshal(body,&netPack)
	if err != nil {
		return &pb.NetRspInfo{
			M_Err:[]byte(err.Error()),
		},nil
	}

	svcAddrs := SubScribeDetail[service]

	if len(svcAddrs)==0{
		return &pb.NetRspInfo{
			M_Err:[]byte(errors.New(fmt.Sprintf("the %s service not found subscribe addresses",service)).Error()),
		},nil
	}

	collectResult := map[uint32]*pb.SendResultInfo{}

	singleResult := make(chan pb.SendResultInfo,len(netPack.M_Net_Pack))

	for key,_ := range netPack.M_Net_Pack {
		go CheckAndPublish(key ,netPack.M_Net_Pack[key],clientIP,service,svcAddrs,singleResult)
	}

	close(singleResult)

	/////读取返回值
	for i := 0;i<len(netPack.M_Net_Pack);i++{
		tmpResult := <- singleResult
		collectResult[tmpResult.Key] = &tmpResult
	}

	return &pb.NetRspInfo{
		M_Net_Rsp:collectResult,
	},nil
}

func CheckAndPublish(key uint32,netpack *pb.Net_Pack,clientIP,service string,svcAddrs []string,result chan pb.SendResultInfo){
	tSendResult := pb.SendResultInfo{
		Key:key,
		CheckErr:nil,
		Result:nil,
	}
	///check ASK_TYPE
	if netpack.M_MsgBody.MLAsktype > uint64(model.ETN_SERVER_SUBSRCTIBE_MSG) {
		tSendResult.CheckErr = []byte(model.ErrAskType.Error())
		result <- tSendResult
		return
	}
	//// check Msg——type
	if netpack.M_MsgBody.MCMsgType >= int32(model.MSG_TYPEMAX){
		tSendResult.CheckErr = []byte(model.ErrMsgType.Error())
		result <- tSendResult
		return
	}
	//// check Send count
	if netpack.M_MsgBody.MSSendCount < 1 {
		tSendResult.CheckErr = []byte(model.ErrSendCount.Error())
		result <- tSendResult
		return
	}
	// 失败了，是否要丢弃
	isDiscard := false
	if netpack.M_MsgBody.MIDiscard == 0 {
		isDiscard = true
	}

	///// 超时时间
	timeout :=  time.Second * time.Duration(netpack.M_MsgBody.MLExpireTime)
	wait := sync.WaitGroup{}
	callResult := make(chan model.BroadcastReturnInfo,len(svcAddrs))

	for _ , tIP := range svcAddrs {
		address := strings.Split(tIP,":")[0]
		port := strings.Split(tIP,":")[1]
		///////=====================
		sendInfo := model.CallInfo{
			ClientIP:clientIP,
			Address:address,
			Port:port,
			Service:service,
			MsgBody:netpack,
			Timeout:timeout,
			IsDiscard:isDiscard,
			AskSequence:netpack.M_MsgBody.MLAskSequence,
			SendTimeApp:netpack.M_MsgBody.MISendTimeApp,
			MsgType :netpack.M_MsgBody.MCMsgType,
			MsgAckType :netpack.M_MsgBody.MCMsgAckType,
			SyncType:model.CALL_CLIENT_PUBLISH,
		}
		wait.Add(1)
		go call.BroadCastClient(sendInfo,callResult,&wait)
	}

	wait.Wait()

	for _,_ = range svcAddrs{
		<- callResult
	}

	tSendResult.CheckErr = nil
	tSendResult.Result = nil
	result <- tSendResult
}

func AsyncAnswer(key uint32,netpack *pb.Net_Pack,suuid string,syncType model.CALL_CLIENT_TYPE,clientIP string,result chan pb.SendResultInfo){

	tSendResult := pb.SendResultInfo{
		Key:key,
		CheckErr:nil,
		Result:nil,
	}
	///check ASK_TYPE
	if netpack.M_MsgBody.MLAsktype > uint64(model.ETN_SERVER_SUBSRCTIBE_MSG) {
		tSendResult.CheckErr = []byte(model.ErrAskType.Error())
		result <- tSendResult
		return
	}
	//// check Msg——type
	if netpack.M_MsgBody.MCMsgType >= int32(model.MSG_TYPEMAX){
		tSendResult.CheckErr = []byte(model.ErrMsgType.Error())
		result <- tSendResult
		return
	}
	//// check Send count
	if netpack.M_MsgBody.MSSendCount < 1 {
		tSendResult.CheckErr = []byte(model.ErrSendCount.Error())
		result <- tSendResult
		return
	}
	// 失败了，是否要丢弃
	isDiscard := false
	if netpack.M_MsgBody.MIDiscard == 0 {
		isDiscard = true
	}

	timeout :=  time.Second * time.Duration(netpack.M_MsgBody.MLExpireTime)

	sendInfo := model.CallInfo{
		Sequence:suuid,
		ClientIP:clientIP,
		SyncType:syncType,
		MsgBody:netpack,
		Timeout:timeout,
		IsDiscard:isDiscard,
		AskSequence:netpack.M_MsgBody.MLAskSequence,
		SendTimeApp:netpack.M_MsgBody.MISendTimeApp,
		MsgType :netpack.M_MsgBody.MCMsgType,
		MsgAckType :netpack.M_MsgBody.MCMsgAckType,
	}

	resultinfo,err := call.AsyncAnswerClient(sendInfo)
	if err != nil {
		tSendResult.CheckErr = []byte(err.Error())
	}

	tSendResult.CheckErr = nil
	tSendResult.Result = resultinfo.M_Net_Rsp
	result <- tSendResult
	return
}
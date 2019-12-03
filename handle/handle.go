package handle

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogo/protobuf/proto"
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

/****
func (m *MsgHandle)ReloadConfig(ctx context.Context, config *pb.Rload)(*pb.Rload,error){
	out := &pb.Rload{}
	utils.ReloadConfig()
	return out , nil
}
*/
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
	err = SqlClient.InsertFunc(in.Sequence)
	if err != nil{
		return nil,err
	}
	err = SqlClient.InsertFuncList(in.Sequence,fmt.Sprintf("%s:%s",in.Ip,in.Port))
	if err != nil{
		return nil,err
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
	err = SqlClient.InsertSvc(in.Service)
	if err != nil {
		return nil,err
	}
	return &pb.PublishReturnInfo{
		Success:true,
	},nil
}

func (m *MsgHandle)Subscribe(ctx context.Context, in *pb.SubscribeInfo)(*pb.SubscribeReturnInfo,error) {
	ip,err := SqlClient.GetSubScribeByIP(in.Service,fmt.Sprintf("%s:%s",in.Ip,in.Port))
	if err != nil{
		return nil,err
	}
	if ip != ""{
		return &pb.SubscribeReturnInfo{
			Success:true,
		},errors.New(fmt.Sprintf("the %s service and %s ip had registered",in.Service,fmt.Sprintf("%s:%s",in.Ip,in.Port)))
	}
	err = SqlClient.InsertSubScribe(in.Service,fmt.Sprintf("%s:%s",in.Ip,in.Port))
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

func ModifyOrFullHead(inbody []byte, info *model.HeadInfo)[]byte{
	if CheckHaveHead(inbody){ // 包含消息头

	}else{ //不包含消息头，应该添加消息头
		headINfo := model.HeadInfo{
			Tag:"ent2015",
			Version:1000,
			ClientType:int16(model.Mid_Msg_Type),
			HeadLength:32,
			CompressWay:uint8(model.Compression_no),
			Encryption:uint8(model.Encryption_No),
			Sig:0,
			Format:0,
			NetFlag:0,
			Back1:0,
			BufSize:int32(len(inbody)),
			UncompressedSize:int32(len(inbody)),
			Back2:0,
		}

	    inbody = utils.JoinHeadAndBody(headINfo,inbody)

	}
	return inbody
}

func AnzalyBodyHead(inbody []byte) (*model.HeadInfo,error) {
	headinfo := model.HeadInfo{}

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
	////check head length
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
	if bufSize > uncompressiondSize {
		return nil,model.ErrCompresseduncompressedLength
	}

	//////如果压缩，解压然后比较长度
	if compressionWay == uint8(model.Compression_zip){
		unzipBytes:=utils.UnzipBytes(inbody[32:])
		if uncompressiondSize != int32(len(unzipBytes)){
			return nil,model.ErrUNCompressedLength
		}
	}
	headinfo.UncompressedSize = uncompressiondSize
	//预留位     4
	m_lBack2 := bodyHead[:]

	headinfo.Back2 = utils.BytesToInt32(m_lBack2)
	return &headinfo,nil
}

func AnzalyBody(inbody []byte,syncType model.CALL_CLIENT_TYPE,clientIP string) (*pb.NetRspInfo,error) {
	body := inbody[32:]
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
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"net pack key:",key)
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MCMsgAckType:",pack.M_MsgBody.MCMsgAckType)
		////model.MSG_TYPE_
		////fmt.Println("pack.M_MsgBody.MCMsgType:",pack.M_MsgBody.MCMsgType)  ///// 消息类型
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MCMsgType:",pack.M_MsgBody.MCMsgType)
		////fmt.Println("pack.M_MsgBody.MIDiscard:",pack.M_MsgBody.MIDiscard)  ///请求可否丢弃// 0：可丢弃 1：不可丢弃
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MIDiscard:",pack.M_MsgBody.MIDiscard)
		////fmt.Println("pack.M_MsgBody.MISendTimeApp:",pack.M_MsgBody.MISendTimeApp) ////开始请求的本地时间戳
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MISendTimeApp:",pack.M_MsgBody.MISendTimeApp)
		////fmt.Println("pack.M_MsgBody.MLAskSequence:",pack.M_MsgBody.MLAskSequence) ////客户请求序列，客户端维护
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MLAskSequence:",pack.M_MsgBody.MLAskSequence)
		////model.ASK_TYPE
		////fmt.Println("pack.M_MsgBody.MLAsktype:",pack.M_MsgBody.MLAsktype)  /// 服务端请求类型
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MLAsktype:",pack.M_MsgBody.MLAsktype)
		////fmt.Println("pack.M_MsgBody.MLBack:",pack.M_MsgBody.MLBack) /////默认为0
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MLBack:",pack.M_MsgBody.MLBack)
		////fmt.Println("pack.M_MsgBody.MLExpireTime:",pack.M_MsgBody.MLExpireTime)  ////过期时间  0：永不过期 >0:过期时间，以m_iSendTimeApp为基本
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MLExpireTime:",pack.M_MsgBody.MLExpireTime)
		////fmt.Println("pack.M_MsgBody.MLResult:",pack.M_MsgBody.MLResult)  /////0：成功 非0：失败
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MLResult:",pack.M_MsgBody.MLResult)
		////fmt.Println("pack.M_MsgBody.MLServerSequence:",pack.M_MsgBody.MLServerSequence) ////服务响应序列(预留)
		//log.TraceWithFields(map[string]interface{}{"func":"AnzalyBody"},"pack.M_MsgBody.MSSendCount:",pack.M_MsgBody.MSSendCount)
		//fmt.Println("pack.M_MsgBody.MSSendCount:",pack.M_MsgBody.MSSendCount)  //// 同一请求次数
		go CheckAndSend(key ,netPack.M_Net_Pack[key],syncType,clientIP,singleResult)
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

func CheckAndSend(key uint32,netpack *pb.Net_Pack,syncType model.CALL_CLIENT_TYPE,clientIP string,result chan pb.SendResultInfo){
	tSendResult := pb.SendResultInfo{
		Key:key,
		SendCount:netpack.M_MsgBody.MSSendCount,
		SuccessCount:0,
		FailCount:0,
		DiscardCount:0,
		ReSendCount:0,
		ResultList:nil,
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

	sendBytes,err := proto.Marshal(netpack)
	if err != nil {
		tSendResult.CheckErr = []byte(err.Error())
		result <- tSendResult
		return
	}
	///// 超时时间
	timeout :=  time.Second * time.Duration(netpack.M_MsgBody.MLExpireTime)

	sendInfo := model.CallInfo{
		ClientIP:clientIP,
		Address:address,
		Port:port,
		Service:service,
		MsgBody:sendBytes,
		Timeout:timeout,
		IsDiscard:isDiscard,
		AskSequence:netpack.M_MsgBody.MLAskSequence,
		SendTimeApp:netpack.M_MsgBody.MISendTimeApp,
		MsgType :netpack.M_MsgBody.MCMsgType,
		MsgAckType :netpack.M_MsgBody.MCMsgAckType,
		SyncType:syncType,
	}

	callResult := make(chan pb.SingleResultInfo,netpack.M_MsgBody.MSSendCount)
	wait := sync.WaitGroup{}
	for i  := 0 ; int32(i) < netpack.M_MsgBody.MSSendCount ; i++{
		wait.Add(1)
		go call.CallClient(sendInfo,callResult,&wait)
	}
	wait.Wait()

	resultList := map[uint32]*pb.SingleResultInfo{}
	failedCount := int32(0)
	discardCount := int32(0)
	resentCount := int32(0)

	for i := 0 ; int32(i) < netpack.M_MsgBody.MSSendCount ; i++{
		tmpRsult := <- callResult
		resultList[uint32(i)] = &tmpRsult
		if tmpRsult.Errinfo != nil {
			failedCount = failedCount + 1
		}
		if tmpRsult.IsDisCard == true {
			discardCount =  discardCount + 1
		}
		if tmpRsult.IsResend == true {
			resentCount = resentCount + 1
		}
	}

	tSendResult.FailCount  		= failedCount
	tSendResult.SuccessCount 	= netpack.M_MsgBody.MSSendCount - failedCount
	tSendResult.ReSendCount 	= resentCount
	tSendResult.DiscardCount 	= discardCount
	tSendResult.ResultList      = resultList

	result <- tSendResult
}

func PublishBody(inbody []byte,service,clientIP string) (*pb.NetRspInfo,error) {
	body := inbody[32:]
	netPack := pb.GJ_Net_Pack{}
	err :=  proto.Unmarshal(body,&netPack)
	if err != nil {
		return &pb.NetRspInfo{
			M_Err:[]byte(err.Error()),
		},nil
	}
	/******
	/////// Get Subscribe by service
	svcAddrs := utils.GetSubscribeByKey(service)
	if len(svcAddrs)==0 {
		log.ErrorWithFields(map[string]interface{}{"func":"PublishBody"},"get services error,not address got")
		return &pb.NetRspInfo{
			M_Err:[]byte(model.ErrGotService.Error()),
		},nil
	}
	*/
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
		SendCount:netpack.M_MsgBody.MSSendCount,
		SuccessCount:0,
		FailCount:0,
		DiscardCount:0,
		ReSendCount:0,
		ResultList:nil,
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

	sendBytes,err := proto.Marshal(netpack)
	if err != nil {
		tSendResult.CheckErr = []byte(err.Error())
		result <- tSendResult
		return
	}

	///// 超时时间
	timeout :=  time.Second * time.Duration(netpack.M_MsgBody.MLExpireTime)
	wait := sync.WaitGroup{}
	callResult := make(chan pb.SingleResultInfo,len(svcAddrs))

	for _ , tIP := range svcAddrs {
		address := strings.Split(tIP,":")[0]
		port := strings.Split(tIP,":")[1]
		///////=====================
		sendInfo := model.CallInfo{
			ClientIP:clientIP,
			Address:address,
			Port:port,
			Service:service,
			MsgBody:sendBytes,
			Timeout:timeout,
			IsDiscard:isDiscard,
			AskSequence:netpack.M_MsgBody.MLAskSequence,
			SendTimeApp:netpack.M_MsgBody.MISendTimeApp,
			MsgType :netpack.M_MsgBody.MCMsgType,
			MsgAckType :netpack.M_MsgBody.MCMsgAckType,
			SyncType:model.CALL_CLIENT_PUBLISH,
		}
		//for i  := 0 ; int32(i) < netpack.M_MsgBody.MSSendCount ; i++{
		wait.Add(1)
		go call.CallClient(sendInfo,callResult,&wait)
		//}
	}
	wait.Wait()


	resultList := make([]pb.SingleResultInfo,0)
	failedCount := int32(0)
	discardCount := int32(0)
	resentCount := int32(0)
	for tmpRsult := range callResult{
		resultList = append(resultList,tmpRsult)
		if tmpRsult.Errinfo != nil {
			failedCount = failedCount + 1
		}
		if tmpRsult.IsDisCard == true {
			discardCount =  discardCount + 1
		}
		if tmpRsult.IsResend == true {
			resentCount = resentCount + 1
		}
	}
	tSendResult.FailCount  		= failedCount
	tSendResult.SuccessCount 	= tSendResult.SendCount - failedCount
	tSendResult.ReSendCount 	= resentCount
	tSendResult.DiscardCount 	= discardCount

	result <- tSendResult
}

package call

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"midmsg/log"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
	"sync"
	"time"
)

////////// 异步调用客户端的call接口
func AsyncCallClient(callinfo model.CallInfo)(*pb.CallRspInfo,error){
	log.Trace("AsyncCallClient")

	caddr := fmt.Sprintf("%v:%v",callinfo.Address,callinfo.Port)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
			//////如果是不丢弃的，超时请求将缓存在队列中
			TimeoutRequest.PutPoolRequest(callinfo)
		}
		return nil,err
	}
	defer conn.Close()

	c := pb.NewClientServiceClient(conn)
	var ctx context.Context
	var cancel context.CancelFunc
	if callinfo.Timeout > time.Second *0 {
		ctx, cancel = context.WithTimeout(context.Background(), callinfo.Timeout)
		defer cancel()
	}else{
		ctx = context.Background()
	}


	///保持异步数据到map中
	StoreAsyncAnswer(callinfo.Sequence,callinfo)

	r, err := c.Call(ctx,&pb.CallReqInfo{M_Body:callinfo.MsgBody,Uuid:callinfo.Sequence,Service:callinfo.Service,Clientip:callinfo.ClientIP})

	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
			//////如果是不丢弃的，超时请求将缓存在队列中
			TimeoutRequest.PutPoolRequest(callinfo)
		}
		log.Error("======================**************err",err.Error())
		return nil,err
	}

	////////////////////超时处理
	if callinfo.Timeout > time.Second *0 {
		select {
		case <-ctx.Done():
			log.Error(ctx.Err()) // 超时处理

			if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
				//////如果是不丢弃的，超时请求将缓存在队列中
				TimeoutRequest.PutPoolRequest(callinfo)
			}else{
				//////如果是丢弃的，超时后返回队列也将丢弃
				LoadAsyncAnswer(callinfo.Sequence)
			}
			return nil,fmt.Errorf("timeout err:%s",ctx.Err().Error())
		default:
			return r,nil
		}
	}

	return r,nil
}
///////////// 异步处理结果失败后，再发起call
func AsyncAnswerClient(callinfo model.CallInfo)(*pb.CallRspInfo,error){

	///////////////////////////调用call async rsp////////////////////////////////////////////////////////////
	log.Trace("callinfo.ClientIP:",callinfo.ClientIP,"utils.ClientPort:",utils.ClientPort)

	loadCallInfo := LoadAsyncAnswer(callinfo.Sequence)

	if loadCallInfo.ClientIP == "" {
		return nil,fmt.Errorf("return info lost ")
	}

	clientAddr := fmt.Sprintf("%v:%d",loadCallInfo.ClientIP,utils.ClientPort)

	clientconn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		////////////将发送失败的异步请求的处理结果，缓存起来
		go AsyncReturn.PutPoolAsyncReturn(callinfo)

		return nil,err
	}
	defer clientconn.Close()

	client := pb.NewClientServiceClient(clientconn)
	var ctxClient context.Context
	ctxClient = context.Background()
	callresult, err := client.AsyncAnswer(ctxClient,&pb.CallReqInfo{M_Body:callinfo.MsgBody,Uuid:callinfo.Sequence,Clientip:callinfo.ClientIP,Service:callinfo.Service})
	if err != nil {
		log.ErrorWithFields(map[string]interface{}{
			"func":"AsyncCallClient",
		},"======================**************AsyncCallClient Err:",err.Error())
		////////////将发送失败的异步请求的处理结果，缓存起来
		go AsyncReturn.PutPoolAsyncReturn(callinfo)

		return nil , err
	}
	return callresult,nil
}

func CallClient(callinfo model.CallInfo)(*pb.CallRspInfo,error){

	caddr := fmt.Sprintf("%v:%v",callinfo.Address,callinfo.Port)

	log.DebugWithFields(map[string]interface{}{"func":"CallClient"},"call client address:",caddr)

	if callinfo.SyncType == model.CALL_CLIENT_ASYNC { ///// 异步
		return AsyncCallClient(callinfo)
	}

	///////////////////////同步操作===========================

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
			TimeoutRequest.PutPoolRequest(callinfo)
		}
		return nil,err
	}
	defer conn.Close()

	c := pb.NewClientServiceClient(conn)
	var ctx context.Context
	var cancel context.CancelFunc
	if callinfo.Timeout > time.Second *0 {
		ctx, cancel = context.WithTimeout(context.Background(), callinfo.Timeout)
		defer cancel()
	}else{
		ctx = context.Background()
	}
	//////////////////////////////////////////////同步

	r, err := c.Call(ctx,&pb.CallReqInfo{M_Body:callinfo.MsgBody,Clientip:callinfo.ClientIP,Service:callinfo.Service,Uuid:callinfo.Sequence})
	log.Error("======================**************同步")
	if err != nil {
		log.Error("======================**************err",err.Error())
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
			TimeoutRequest.PutPoolRequest(callinfo)
		}

		return nil,err
	}
	////////////////////超时处理
	if callinfo.Timeout > time.Second *0 {
		select {
		case <-ctx.Done():
			if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
				TimeoutRequest.PutPoolRequest(callinfo)
				//////////丢弃了
			}
			return  nil,fmt.Errorf("timeout")
		default:
			return r,nil
		}
	}
	return r,nil
}

func BroadCastClient(callinfo model.CallInfo,rsp chan model.BroadcastReturnInfo,wait *sync.WaitGroup){
	defer wait.Done()
	caddr := fmt.Sprintf("%v:%v",callinfo.Address,callinfo.Port)
	log.DebugWithFields(map[string]interface{}{"func":"BroadCastClient"},"broadcast client address:",caddr)

	callRsp := model.BroadcastReturnInfo{
		SResult:nil,
		SErr:nil,
	}

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
			TimeoutRequest.PutPoolRequest(callinfo)
		}
		callRsp.SErr = err
		rsp <- callRsp
		return
	}
	defer conn.Close()

	c := pb.NewClientServiceClient(conn)
	var ctx context.Context
	var cancel context.CancelFunc
	if callinfo.Timeout > time.Second *0 {
		ctx, cancel = context.WithTimeout(context.Background(), callinfo.Timeout)
		defer cancel()
	}else{
		ctx = context.Background()
	}
	//////////////////////////////////////////////同步

	r, err := c.Call(ctx,&pb.CallReqInfo{M_Body:callinfo.MsgBody,Clientip:callinfo.ClientIP,Service:callinfo.Service,Uuid:callinfo.Sequence})
	log.Error("======================**************广播")
	if err != nil {
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
			TimeoutRequest.PutPoolRequest(callinfo)
		}
		log.Error("======================**************err",err.Error())
		callRsp.SErr = err
		rsp <- callRsp
		return
	}
	////////////////////超时处理
	if callinfo.Timeout > time.Second *0 {
		select {
		case <-ctx.Done():
			if callinfo.IsDiscard != true { ///// 超时了不可丢弃放到 重新发送的pool里
				TimeoutRequest.PutPoolRequest(callinfo)
			}
			callRsp.SErr = fmt.Errorf("request timeout")
			rsp <- callRsp
			return
		default:
			callRsp.SResult = r
			callRsp.SErr = nil
			return
		}
	}
	callRsp.SResult = r
	callRsp.SErr = nil
	return
}

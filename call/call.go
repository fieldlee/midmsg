package call

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
	"sync"
	"time"
)



func AsyncCallClient(callinfo model.CallInfo){
	caddr := fmt.Sprintf("%v:%v",callinfo.Address,callinfo.Port)
	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
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

	sResult := pb.SingleResultInfo{
		AskSequence:callinfo.AskSequence,
		SendTimeApp:callinfo.SendTimeApp,
		MsgType:	callinfo.MsgType,
		MsgAckType:	callinfo.MsgAckType,
		SyncType:	callinfo.SyncType,
		IsTimeOut:false,
		IsDisCard:false,
		IsResend:false,
		Errinfo:nil,
		Result:nil,
	}


	r, err := c.Call(ctx,&pb.NetReqInfo{M_Body:callinfo.MsgBody})

	//////////////////////异步处理 ， 调用客户端的接口，异步发送
	if err != nil {
		sResult.Errinfo = []byte(err.Error())
	}else{
		sResult.Result = r.M_Net_Rsp
	}

	////////////////////超时处理
	if callinfo.Timeout > time.Second *0 {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err()) // 超时处理
			sResult.IsTimeOut = true
			if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
				sResult.IsResend = true
				PutPoolRequest(callinfo)
			}else{
				sResult.IsDisCard = true
			}
		}
	}

	///////////////////////////调用call async rsp////////////////////////////////////////////////////////////
	clientAddr := fmt.Sprintf("%v:%d",callinfo.ClientIP,utils.ClientPort)
	clientconn, err := grpc.Dial(clientAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer clientconn.Close()

	client := pb.NewClientServiceClient(clientconn)
	var ctxClient context.Context
	ctxClient = context.Background()
	_, err = client.AsyncCall(ctxClient,&sResult)
	if err != nil {
		return
	}

}

func CallClient(callinfo model.CallInfo, tResult chan pb.SingleResultInfo, wait *sync.WaitGroup){
	if wait != nil {
		defer wait.Done()
	}

	sResult := pb.SingleResultInfo{
		AskSequence:callinfo.AskSequence,
		SendTimeApp:callinfo.SendTimeApp,
		MsgType:callinfo.MsgType,
		MsgAckType:callinfo.MsgAckType,
		SyncType:callinfo.SyncType,
		IsTimeOut:false,
		IsDisCard:false,
		IsResend:false,
		Errinfo:nil,
		Result:nil,
	}

	caddr := fmt.Sprintf("%v:%v",callinfo.Address,callinfo.Port)
	fmt.Println("call client address:",caddr)
	if sResult.SyncType == 1 { ///// 异步
		/// 异步调用goroutine
		go AsyncCallClient(callinfo)

		if tResult != nil {
			sResult.Errinfo = nil
			tResult <- sResult
		}
		return
	}

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		if tResult != nil {
			sResult.Errinfo = []byte(err.Error())
			tResult <- sResult
		}
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
	r, err := c.Call(ctx,&pb.NetReqInfo{M_Body:callinfo.MsgBody})

	if err != nil {
		if tResult != nil {
			sResult.Errinfo = []byte(err.Error())
			tResult <- sResult
		}
		return
	}else{
		if tResult != nil {
			fmt.Println("===============================================================")
			fmt.Println(string(r.M_Net_Rsp))
			//////// 是否将结果返回到客户端  服务器等  /0 无需回复, 1 回复到发送方, 2 回复到离线服务器
			if sResult.MsgAckType  == 1 {
				sResult.Result = r.M_Net_Rsp
			}
			sResult.Errinfo = nil
			tResult <- sResult
		}
		return
	}
	////////////////////超时处理
	select {
	case <-ctx.Done():

		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
			PutPoolRequest(callinfo)
			//////////丢弃了
			if tResult != nil {
				sResult.IsResend = true
			}
		}else{
			//////////丢弃了
			if tResult != nil {
				sResult.IsDisCard = true
			}
		}

		if tResult != nil {
			sResult.Result = nil
			sResult.IsTimeOut = true
			sResult.Errinfo = []byte(ctx.Err().Error())
			tResult <- sResult
		}

		return
	}
	if tResult != nil {
		sResult.Result = nil
		sResult.IsTimeOut = false
		sResult.Errinfo = nil
		tResult <- sResult
	}

	return
}
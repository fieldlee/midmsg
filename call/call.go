package call

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"midmsg/model"
	pb "midmsg/proto"
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
	_, err = c.Call(ctx,&pb.NetReqInfo{M_Body:callinfo.MsgBody})
	if err != nil {
		return
	}
	////////////////////超时处理
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // 超时处理
		if callinfo.IsDiscard != true { ///// 超时了不可丢弃放在 重新发送的pool里
			PutPoolRequest(callinfo)
		}
		return
	}
}

func CallClient(callinfo model.CallInfo, tResult chan model.SingleResultInfo, wait *sync.WaitGroup){
	if wait != nil {
		defer wait.Done()
	}

	sResult := model.SingleResultInfo{
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
			sResult.Errinfo = err
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
			sResult.Errinfo = err
			tResult <- sResult
		}
		return
	}else{
		if tResult != nil {

			//////// 是否将结果返回到客户端  服务器等  /0 无需回复, 1 回复到发送方, 2 回复到离线服务器

			if sResult.MsgAckType  == 1 {
				sResult.Result = r
			}

			sResult.Errinfo = nil
			tResult <- sResult
		}
		return
	}
	////////////////////超时处理
	select {
	case <-ctx.Done():

		fmt.Println(ctx.Err()) // 超时处理

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
			sResult.Errinfo = ctx.Err()
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
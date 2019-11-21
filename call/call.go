package call

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "midmsg/proto"
	"time"
)


func CallClient(address,port string,timeout time.Duration,service string,msg []byte)([]byte,error){

	caddr := fmt.Sprintf("%v:%v",address,port)

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()

	c := pb.NewClientServiceClient(conn)
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > time.Second *0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}else{
		ctx = context.Background()
	}

	r, err := c.Call(ctx,&pb.NetReqInfo{M_Body:msg})

	if err != nil {
		return nil,err
	}else{
		return r.M_Resp,nil
	}
	////////////////////超时处理
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err()) // 超时处理
		PutPoolRequest(address,port,timeout,service,msg)
	}
	return nil,nil

}
package call

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "midmsg/proto"
	"time"
)

func CallClient(address,port string,timeout time.Duration,msg []byte)([]byte,error){

	caddr := fmt.Sprintf("%v:%v",address,port)

	conn, err := grpc.Dial(caddr, grpc.WithInsecure())
	if err != nil {
		return nil,err
	}
	defer conn.Close()

	c := pb.NewClientServiceClient(conn)
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > time.Second {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}else{
		ctx = context.Background()
	}

	r, err := c.Call(ctx,&pb.NetReqInfo{M_Body:msg})

	if err != nil {
		return nil,err
	}

	return r.M_Resp,nil
}
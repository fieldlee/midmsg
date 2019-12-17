package handle

import (
	"midmsg/log"
	"midmsg/model"
	pb "midmsg/proto"
)

func NewWorker(workerPool chan chan HandleBody, jobDone chan struct{}) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobDone:    jobDone,
		JobChannel: make(chan HandleBody),
		quit:       make(chan bool),
	}
}

func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel
			select {
			case body := <-w.JobChannel:
				// we have received a work request.
				//log.Error("====================GoRoutineId:",utils.GetGID())
				// 解析头文件
				if headInfo,err := AnzalyBodyHead(body.MBody); err != nil {
					log.ErrorWithFields(map[string]interface{}{
						"func":"Worker.start",
					},err.Error())
					pbRespinfo := &pb.NetRspInfo{
						M_Err:[]byte(err.Error()),
					}
					go func(info *pb.NetRspInfo) {
						body.Out <- info
					}(pbRespinfo)

				}else{ ///////// 校验package head 完成后 校验package 内容
					inBody := ModifyOrFullHead(body.MBody,headInfo)  /////修改后的包bytes

					if body.Type == model.CALL_CLIENT_PUBLISH {
						/// 订阅消息发送
						rspInfo,err := PublishBody(inBody,body.Service,body.ClientIp)
						if err != nil {
							pbRespinfo := &pb.NetRspInfo{
								M_Err:[]byte(err.Error()),
							}
							go func(info *pb.NetRspInfo) {
								body.Out <- info
							}(pbRespinfo)
						}
						go func(info *pb.NetRspInfo) {
							body.Out <- info
						}(rspInfo)
					}else{
						//////异步回复和同步调用
						rspInfo,err := AnzalyBody(inBody,body.Sequence,body.Type,body.ClientIp)
						if err != nil {
							pbRespinfo := &pb.NetRspInfo{
								M_Err:[]byte(err.Error()),
							}
							go func(info *pb.NetRspInfo) {
								body.Out <- info
							}(pbRespinfo)
						}
						go func(info *pb.NetRspInfo) {
							body.Out <- info
						}(rspInfo)
					}
				}

				w.JobDone <- struct{}{}

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
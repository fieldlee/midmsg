package handle

import (
	"fmt"
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
				// 解析头文件
				err := AnzalyBodyHead(body.MBody)
				if err != nil {
					fmt.Println(err.Error())
					pbRespinfo := &pb.NetRspInfo{
						M_Err:[]byte(err.Error()),
					}
					go func(info *pb.NetRspInfo) {
						body.Out <- info
					}(pbRespinfo)
				}else{ ///////// 校验package head 完成后 校验package 内容
					rspInfo,err := AnzalyBody(body.MBody,uint32(body.Type),body.ClientIp)
					if err != nil {
						pbRespinfo := &pb.NetRspInfo{
							M_Err:[]byte(err.Error()),
						}
						go func(info *pb.NetRspInfo) {
							body.Out <- info
						}(pbRespinfo)
						//return
					}
					go func(info *pb.NetRspInfo) {
						body.Out <- info
					}(rspInfo)
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
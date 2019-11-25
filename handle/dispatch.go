package handle

import (
	"midmsg/log"
	"midmsg/model"
	pb "midmsg/proto"
	"midmsg/utils"
)

type HandleBody struct {
	MBody 		[]byte
	Type   		model.CALL_CLIENT_TYPE
	ClientIp 	string
	Out    		chan *pb.NetRspInfo
}

type Dispatcher struct {
	WorkerPool  chan chan HandleBody
	MaxWork 	uint32
	JobDone     chan struct{}
}

type Worker struct {
	WorkerPool  chan chan HandleBody
	JobChannel  chan HandleBody
	JobDone     chan struct{}
	quit    	chan bool
}
// A buffered channel that we can send work requests on.
var (
	JobQueue = make(chan HandleBody, utils.MaxQueue)
 	JobDone = make(chan struct{}, utils.MaxWorker)
)

func NewDispatcher(maxWorkers uint32,jobDone chan struct{}) *Dispatcher {
	pool := make(chan chan HandleBody, maxWorkers)
	return &Dispatcher{WorkerPool:pool,MaxWork:maxWorkers,JobDone:jobDone}
}

func (d *Dispatcher) Run() {
	log.Info("Worker queue dispatcher started...")
	// starting n number of workers
	for i := 0; i < int(d.MaxWork); i++ {
		worker := NewWorker(d.WorkerPool,d.JobDone)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			// a job request has been received
			go func(job HandleBody) {
				jobChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}

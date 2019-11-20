package handle

import (
	"midmsg/model"
	pb "midmsg/proto"
)

type HandleBody struct {
	M_Body []byte
	Type   model.CALL_CLIENT_TYPE
	Out    chan *pb.NetRspInfo
	Err    chan error
}

type Dispatcher struct {
	WorkerPool chan chan HandleBody
	MaxWork uint32
}

type Worker struct {
	WorkerPool  chan chan HandleBody
	JobChannel  chan HandleBody
	quit    	chan bool
}
// A buffered channel that we can send work requests on.
var JobQueue chan HandleBody

func NewDispatcher(maxWorkers uint32) *Dispatcher {
	pool := make(chan chan HandleBody, maxWorkers)
	return &Dispatcher{WorkerPool:pool,MaxWork:maxWorkers}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < int(d.MaxWork); i++ {
		worker := NewWorker(d.WorkerPool)
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
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job

			}(job)
		}
	}
}
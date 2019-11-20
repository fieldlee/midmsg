package handle

func NewWorker(workerPool chan chan HandleBody) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan HandleBody),
		quit:       make(chan bool)}
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
				err := AnzalyBodyHead(body.M_Body)
				if err != nil {
					body.Err <- err
				}


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

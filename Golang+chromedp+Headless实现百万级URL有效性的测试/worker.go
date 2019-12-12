package main

type Worker struct {
    JobQueue chan Job
}

func NewWorker() Worker {
    return Worker{JobQueue: make(chan Job)}
}
func (w Worker) Run(wq chan chan Job) {
    go func() {
        // defer wg.Add(-1)
        for {
            wq <- w.JobQueue
            select {
            case job := <-w.JobQueue:
                job.Do()
                // wg.Add(-1)
            }
        }
    }()
}
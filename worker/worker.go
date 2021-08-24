package worker

import (
    "fmt"
    "math/rand"
    "time"
)

type Worker struct {
    id  int
    cid int
}

func New(myID, controllerID int) *Worker {
    fmt.Printf("New worker %d\n", myID)
    return &Worker {
        id: myID,
        cid: controllerID,
    }
}

func (w *Worker) Run(done <-chan struct{}, workItems <-chan string, results chan<- error) {
    for workItem := range workItems {
        fmt.Printf("Worker %d:%d got a workitem: %s\n", w.cid, w.id, workItem)
        select {
        case results <- w.work(workItem):
            fmt.Printf("Result sent from %d:%d\n", w.cid, w.id)
        case <- done:
            fmt.Printf("Worker %d:%d exiting on done signal\n", w.cid, w.id)
            return
        }
    }
    fmt.Printf("Worker %d:%d exiting, no more work items.\n", w.cid, w.id)
}

func (w *Worker) work(item string) error {
    sleep := time.Duration(rand.Intn(10)) * time.Second
    fmt.Printf("Worker %d:%d processing %s [%v]\n", w.cid, w.id, item, sleep)
    time.Sleep(sleep)

    fmt.Printf("Worker %d:%d finishing\n", w.cid, w.id)
    if rand.Intn(4) == 0 {
        return fmt.Errorf("worker %d:%d encountered an error", w.cid, w.id)
    }
    return nil
}

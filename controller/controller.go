package controller

import (
    "context"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    "github.com/bnaylor/fanout/worker"
)

type Controller struct {
    id       int
}

func New(myID int) *Controller {
    return &Controller{
        id: myID,
    }
}

func (c *Controller) ID() int {
    return c.id
}

func (c *Controller) Run(ctx context.Context, poolCount int, root string) (int, error) {
    results := make(chan error)

    // Start up the first stage (collect some filenames)
    paths, walkErrors := walkFiles(ctx, root)

    // Second stage, fan out a bunch of workers
    var wg sync.WaitGroup
    wg.Add(poolCount)

    for i := 0; i < poolCount; i++ {
        go func(id int) {
            w := worker.New(id, c.id)
            w.Run(ctx, paths, results)
            wg.Done()
        }(i)
    }

    go func() {
        fmt.Printf("Controller %d waiting for workers.\n", c.id)
        wg.Wait()
        fmt.Printf("Controller %d done waiting.\n", c.id)
        close(results)
    }()

    // Consume the results
    fmt.Printf("Controller %d processing results..\n", c.id)
    resultCount := 0
    for result := range results {
        msg := "OK."
        if result != nil {
            msg = result.Error()
        }
        fmt.Printf("Result: %s\n", msg)
        resultCount++
    }

    // Note: we can only consume the walk errors at this point
    fmt.Printf("Checking for walk errors.\n")
    if err := <- walkErrors; err != nil {
        return resultCount, fmt.Errorf("there was a walk error: %v", err)
    }

    return resultCount, nil
}

// walkFiles gets us some strings to "process"
func walkFiles(ctx context.Context, root string) (<-chan string, <-chan error) {
    paths := make(chan string)
    errc := make(chan error, 1)

    go func() {
        // Close the paths channel after Walk returns.
        defer walkFinished(paths)

        // No select needed for this send, since errc is buffered.
        errc <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.Mode().IsRegular() {
                return nil
            }
            select {
            case paths <- path:
            case <-ctx.Done():
                return errors.New("walk canceled")
            }
            return nil
        })
    }()
    return paths, errc
}

func walkFinished(c chan string) {
    fmt.Printf("Walk finished, closing channel\n")
    close(c)
}

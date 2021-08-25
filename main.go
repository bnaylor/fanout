package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/bnaylor/fanout/controller"
)

const fileRoot = "/Users/bnaylor/src/detritus/docker"

var (
	controllers    []*controller.Controller
	numControllers = 2
	cancelFunc     context.CancelFunc
)

func main() {
	// Seed RNG one time
	rand.Seed(time.Now().UnixNano())

	// Signal handlers
	installHandlers()

	ctx := context.Background()
	ctx, cancelFunc = context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(numControllers)

	counts := make(chan int)

	// Go
	controllers = make([]*controller.Controller, 0)
	for i := 0; i < numControllers; i++ {
		cont := controller.New(i)
		controllers = append(controllers, cont)
		go func(c *controller.Controller) {
			fmt.Printf("Starting controller %d..\n", c.ID())
			processed, err := c.Run(ctx, 16, fileRoot)
			fmt.Printf("%d results from controller %d\n", processed, c.ID())
			counts <- processed
			if err != nil {
				fmt.Printf("Error running controller %d: %v\n", c.ID(), err)
				cancelFunc()
			}
			wg.Done()
		}(cont)
	}

	defer func() {
		wg.Wait()
		close(counts)
	}()

	total := 0
	for i := 0; i < numControllers; i++ {
		total += <-counts
	}
	fmt.Printf("Done, %d files processed.\n", total)
}

func installHandlers() {
	go sigQuitHandler()
}

func sigQuitHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT)

	for {
		_ = <-sigs
		fmt.Printf("sigQ\n")
		dumpGoroutines()

		// This is not actually necessary, just seeing how it works
		fmt.Printf("Shutting down %d controllers \n", len(controllers))
		for _, c := range controllers {
			fmt.Printf("Sending cancel to controller %d.\n", c.ID())
			cancelFunc()
		}
	}
}

func dumpGoroutines() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	fmt.Printf("\n=== Begin Goroutine Dump ===\n%s\n=== End Goroutine Dump ===\n", string(buf[:stacklen]))
}

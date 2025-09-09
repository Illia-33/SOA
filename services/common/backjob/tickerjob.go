package backjob

import (
	"context"
	"fmt"
	"log"
	"time"
)

type JobCallback func(context.Context) error

type TickerJob struct {
	ticker   *time.Ticker
	job      JobCallback
	commands chan command
}

func NewTickerJob(period time.Duration, j JobCallback) TickerJob {
	return TickerJob{
		ticker:   time.NewTicker(period),
		job:      j,
		commands: make(chan command),
	}
}

func (w *TickerJob) Run() {
	go func() {
		w.run()
	}()
}

func (w *TickerJob) Stop() {
	w.commands <- command_stop
	close(w.commands)
	w.ticker.Stop()
}

func (w *TickerJob) ForceCheckOutbox() {
	w.commands <- command_force_do_job
}

func (w *TickerJob) work() (bool, error) {
	select {
	case command := <-w.commands:
		{
			switch command {
			case command_stop:
				{
					return false, nil
				}

			case command_force_do_job:
				{
					return true, w.job(context.Background())
				}
			}

			return true, fmt.Errorf("unknown command: %d", int(command))
		}

	case <-w.ticker.C:
		{
			return true, w.job(context.Background())
		}
	}
}

func (w *TickerJob) run() {
	for {
		continueFlag, err := w.work()

		if err != nil {
			log.Printf("background job error: %v", err)
		}

		if !continueFlag {
			log.Println("background job stopped")
			break
		}
	}
}

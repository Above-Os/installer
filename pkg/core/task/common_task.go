package task

import (
	"context"
	"fmt"
	"time"

	"bytetrade.io/web3os/installer/pkg/core/action"
)

type CommonTask struct {
	Name    string
	Desc    string
	Action  action.Action
	Retry   int
	Delay   time.Duration
	Timeout time.Duration
}

func (t *CommonTask) GetDesc() string {
	return t.Desc
}

func (t *CommonTask) Init() {
	t.Default()
}

func (t *CommonTask) Default() {

}

func (t *CommonTask) Execute() {

}

func (t *CommonTask) RunWithTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	resCh := make(chan error)

	go t.Run(resCh)
	select {
	case <-ctx.Done():
		// t.TaskResult.AppendErr(host, fmt.Errorf("execute task timeout, Timeout=%s", util.ShortDur(l.Timeout)))
	case e := <-resCh:
		fmt.Println("---e---", e)
		// if e != nil {
		// 	t.TaskResult.AppendErr(host, e)
		// }
	}
}

func (t *CommonTask) Run(resCh chan error) {
	var res error
	defer func() {
		resCh <- res
		close(resCh)
	}()
}

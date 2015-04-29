package workflow

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
)

type Duty map[string]string
type Workflow struct {
	Context context.Context
	Queue   map[string]Duty
	Now     string
	Deep    uint32
	once    sync.Once
}

func New(ctx context.Context) Workflow {
	wf := Workflow{}
	wf.Context = ctx
	wf.Queue = make(map[string]Duty)
	return wf
}

func (wf *Workflow) Add(name string, duty Duty) {
	wf.Queue[name] = duty
}

func (wf *Workflow) Handle(status string) error {
	wf.Deep = 0
	wf.Now = status
	return wf.next()
}
func (wf *Workflow) next() error {
	for {
		atomic.AddUint32(&wf.Deep, 1)
		// agent process
		if wf.Now == "end" {
			return nil
		}
		agent, ok := GetAgent(wf.Now)
		if !ok {
			return errors.New("agent named [" + wf.Now + "] not found")
		}
		var (
			next string
			err  error
		)
		wf.Context, next, err = agent(wf.Context)
		log.Println(wf.Now, next)
		if err != nil {
			return err
		}

		// find next step
		d, ok := wf.Queue[wf.Now]
		if !ok {
			return nil
		}
		wf.Now, _ = d[next]
	}

	return nil

}

package main

import (
	"log"
	"net/http"
	"os"
	"thumbnail"
	"time"
	"workflow"

	"golang.org/x/net/context"
)

type MuxContext struct {
	context.Context
	Cancel func()
}

func (p *MuxContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	statusCode := http.StatusOK
	var err interface{}
	defer func() {
		endTime := time.Now()
		if err == nil {
			err = ""
		}
		log.Printf("%d %0.3f %s %s [%s] [%s]", statusCode, float32(endTime.Sub(startTime))/float32(time.Second), r.Method, r.URL.Path, r.URL.RawQuery, err)
	}()

	ctx := context.WithValue(context.WithValue(context.Background(), "req", r), "res", w)
	wf := workflow.New(ctx)
	if err != nil {
		return
	}
	registWorkflow(wf)
	err = wf.Handle("start")
}
func getMux() *MuxContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &MuxContext{ctx, cancel}
}

func registWorkflow(wf workflow.Workflow) {
	var d = workflow.Duty{}
	d["ok"] = "thumbnail"
	wf.Add("start", d)
	d = workflow.Duty{}
}

func agentMount() {
	workflow.NewAgent(thumbnail.Thumbnail, "thumbnail")
	workflow.NewAgent(func(ctx context.Context) (context.Context, string, error) {
		return ctx, "ok", nil
	}, "start")
}

func main() {
	var port = "8080"

	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	// start server
	ctx := getMux()
	agentMount()
	go func() {
		err := http.ListenAndServe(":"+port, ctx) //设置监听的端口
		if err != nil {
			ctx.Cancel()
			log.Print(err)
		}
	}()
	log.Println("service start at ", port)
	<-ctx.Done()
	// http://timg.baidu.com/timg?sec=1430287317&di=acf0a0802b9c6d5da2b7bee9a0adaa0f&size=f200_200&quality=90&src=http%3A%2F%2Fstatic.qyy.baidu.com%2Fstatic%2Ffang%2Fhome%2Fstatic%2Fapp%2F5753624%2Fhuxingtu%2F2_470edf9.jpg
}

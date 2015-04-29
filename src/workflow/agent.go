package workflow

import "golang.org/x/net/context"

var (
	agentPool = make(map[string]Agent)
)

type Agent func(context.Context) (context.Context, string, error)

func NewAgent(fn Agent, alias ...string) {
	for _, name := range alias {
		agentPool[name] = fn
	}
}
func GetAgent(name string) (Agent, bool) {
	if ag, ok := agentPool[name]; ok {
		return ag, true
	}
	return nil, false
}
